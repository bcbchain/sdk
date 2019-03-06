package dockerlib

import (
	"net"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/tendermint/tmlibs/log"
	"golang.org/x/net/context"
)

// TODO: 目前達爾文是 linux 的拷貝，水果系統上未見與 linux 不同的地方，記得一起改

// DockerLib 是我們自定義的 Docker API 的 Wrapper
type DockerLib struct {
	logger log.Logger
}

var (
	myLib        *DockerLib
	instanceOnce sync.Once
	initOnce     sync.Once
)

// GetDockerLib 初始化得到 DockerLib 對象指針
func GetDockerLib() *DockerLib {
	instanceOnce.Do(func() {
		myLib = &DockerLib{}
	})
	return myLib
}

// Init 傳入日志對象，不能不傳
func (l *DockerLib) Init(log log.Logger) {
	initOnce.Do(func() {
		l.logger = log
	})
}

// GetMyIntranetIP 獲得本機局網網卡 IP，如有多個，取第一個
func (l *DockerLib) GetMyIntranetIP() string {
	addrArray, err := net.InterfaceAddrs()
	if err != nil {
		l.logger.Warn("DockerLib GetMyIntranetIP cause ERROR", "err", err)
		return ""
	}
	for _, addr := range addrArray {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}

// GetDockerHubIP 獲得本機 Docker 的網卡地址，如果有服務需要 Docker 容器內部訪問，就可以訪問這個地址
func (l *DockerLib) GetDockerHubIP() string {
	params := DockerRunParams{
		Cmd:        []string{"ip", "r"},
		NeedOut:    true,
		NeedWait:   true,
		NeedRemove: true,
	}
	if !l.Run("alpine:latest", "", &params) {
		return ""
	}
	listStr := strings.Split(params.FirstOutput, " ")
	if len(listStr) < 5 {
		l.logger.Warn("GetDockerHubIP got strange output:", "stdout", params.FirstOutput)
	}
	return listStr[2] // this is the result
}

// Mounts DockerRun所需目錄映射
type Mounts struct {
	Source      string
	Destination string
	ReadOnly    bool
}

// HostPort DockerRun 需要映射到的本機 IP 和 端口
type HostPort struct {
	Port string
	Host string
}

// DockerRunParams 運行 Docker 容器需要的參數，避免調用者還依賴 Docker API
type DockerRunParams struct {
	Env         []string
	Cmd         []string
	WorkDir     string
	Mounts      []Mounts
	PortMap     map[string]HostPort
	FirstOutput string // 回寫
	NeedOut     bool   // 需要拿到控制臺輸出（只拿開始的內容，不能一直等，有些進程會一直輸出）
	NeedRemove  bool   // 需要手工清理掉屍體
	AutoRemove  bool   // 給 daemon 設置一下，如果它們掛了，就自己打掃戰場，不留垃圾
	NeedWait    bool   // 等它執行結束（需要注意 daemon 不會結束）
}

// Run 運行 Docker 容器，執行某個功能。由於無法直接獲知Docker內Service的啓動狀態，請參考test文件中的處理辦法，或者在Service啓動的時候主動回調
func (l *DockerLib) Run(dockerImageName, containerName string, params *DockerRunParams) bool {
	l.logger.Info("DockerLib Run", "image", dockerImageName, "containerName", containerName, "params", params)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		l.logger.Warn("DockerLib Run NewEnvClient Error:", "err", err)
		return false
	}

	if !l.ensureImage(ctx, cli, dockerImageName) {
		return false
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        dockerImageName,
			Cmd:          params.Cmd,
			Tty:          true,
			Env:          params.Env,
			WorkingDir:   params.WorkDir,
			ExposedPorts: assemblePortSet(params),
		}, &container.HostConfig{
			Mounts:       assembleMounts(params),
			PortBindings: assemblePortMap(params),
			AutoRemove:   params.AutoRemove,
		}, nil, containerName)
	if err != nil {
		l.logger.Warn("DockerLib Run ContainerCreate Error:", "err", err)
		return false
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		l.logger.Warn("DockerLib Run ContainerStart Error:", "err", err)
		return false
	}

	if params.NeedWait {
		if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
			l.logger.Warn("DockerLib Run ContainerWait Error:", "err", err)
			return false
		}
	}

	if !l.feedBack(ctx, cli, resp.ID, params) {
		return false
	}

	if params.NeedRemove {
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
		if err != nil {
			l.logger.Warn("DockerLib Run remove cause ERROR:", "err", err, "Please remove manually", containerName)
		}
	}

	return true
}

func (l *DockerLib) feedBack(ctx context.Context, cli *client.Client, containerID string, params *DockerRunParams) bool {
	if params.NeedOut {
		out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			l.logger.Warn("DockerLib Run ContainerLogs cause ERROR:", "err", err)
			return false
		}

		byt := make([]byte, 3000)
		n, err := out.Read(byt)
		if err != nil {
			l.logger.Warn("DockerLib Run Read From ContainerLogs cause ERROR:", "err", err)
		}
		if n <= 0 {
			l.logger.Warn("DockerLib Run Read From ContainerLogs cause ERROR: output is zero length")
		}
		params.FirstOutput = string(byt)
	}
	return true
}

func assemblePortSet(params *DockerRunParams) nat.PortSet {
	portSet := make(map[nat.Port]struct{}, 0)
	for k := range params.PortMap {
		p := nat.Port(k)
		portSet[p] = struct{}{}
	}
	return portSet
}

func assemblePortMap(params *DockerRunParams) nat.PortMap {
	portMap := make(map[nat.Port][]nat.PortBinding, 0)
	for k, v := range params.PortMap {
		p := nat.Port(k)
		bindings := make([]nat.PortBinding, 1)
		bindings[0] = nat.PortBinding{
			HostIP:   v.Host,
			HostPort: v.Port,
		}
		portMap[p] = bindings
	}
	return portMap
}

func assembleMounts(params *DockerRunParams) []mount.Mount {
	mounts := make([]mount.Mount, 0)
	for _, m := range params.Mounts {
		mt := mount.Mount{Type: mount.TypeBind,
			Source:   m.Source,
			Target:   m.Destination,
			ReadOnly: m.ReadOnly,
		}
		mounts = append(mounts, mt)
	}
	return mounts
}

func (l *DockerLib) ensureImage(ctx context.Context, cli *client.Client, imageName string) bool {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		l.logger.Warn("DockerLib Run ImageList Error:", "err", err)
		return false
	}

	if notExists(images, imageName) {
		p, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		defer p.Close()
		if err != nil {
			l.logger.Warn("DockerLib Run ImagePull Error:", "err", err)
			return false
		}

		byt := make([]byte, 500)
		for {
			_, err := p.Read(byt)
			if err != nil {
				l.logger.Info("DockerLib ImagePull can't Read output", "err", err)
				break
			} else {
				if strings.Contains(string(byt), "Downloaded") || strings.Contains(string(byt), "up to date") {
					l.logger.Debug("DockerLib ImagePull:", "result", string(byt))
					break
				}
			}
		}
	}
	return true
}

func notExists(images []types.ImageSummary, imageName string) bool {
	exists := false
	for _, image := range images {
		// fmt.Println(image.RepoTags)
		for _, tag := range image.RepoTags {
			if tag == imageName {
				exists = true
				break
			}
		}
		if exists {
			break
		}
	}
	return !exists
}

// Kill 殺死一個 Docker 容器，並且清理現場
func (l *DockerLib) Kill(containerName string) bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		l.logger.Warn("DockerLib Kill NewEnvClient cause ERROR:", "err", err)
		return false
	}

	containerID := l.getContainerIDByName(ctx, cli, containerName)
	if containerID == "" {
		l.logger.Warn("No such containerName:", "name", containerName)
		return true // 木有的情況也返回 true 吧，就省了 remove 了
	}

	return l.killByID(ctx, cli, containerID)
}

func (l *DockerLib) killByID(ctx context.Context, cli *client.Client, containerID string) bool {
	err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		l.logger.Warn("DockerLib Kill remove cause ERROR:", "err", err, "Please remove manually", containerID)
		return false
	}
	return true
}

func (l *DockerLib) getContainerIDByName(ctx context.Context, cli *client.Client, containerName string) string {
	containerID := ""
	list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		l.logger.Warn("DockerLib getContainerIDByName cause ERROR:", "err", err)
		return ""
	}
	for _, con := range list {
		for _, name := range con.Names {
			if name[1:] == containerName {
				containerID = con.ID
				break
			}
		}
		if containerID != "" {
			break
		}
	}
	return containerID
}

// Status 查詢一個容器的狀態
func (l *DockerLib) Status(containerName string) bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		l.logger.Warn("DockerLib Status NewEnvClient cause ERROR:", "err", err)
		return false
	}

	containerID := l.getContainerIDByName(ctx, cli, containerName)
	if containerID == "" {
		l.logger.Warn("No such containerName:", "name", containerName)
		return false
	}
	stat, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		l.logger.Warn("DockerLib Status ContainerStats cause ERROR:", "err", err)
		return false
	}
	if stat.Body == nil {
		return false
	}
	err = stat.Body.Close()
	if err != nil {
		l.logger.Debug("DockerLib Status Close ContainerStats response:", "err", err)
	}

	return true
}

// Reset 殺掉所有自己啓動的容器(以特定字冠命名的)
func (l *DockerLib) Reset(prefix string) bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		l.logger.Warn("DockerLib Reset NewEnvClient cause ERROR:", "err", err)
		return false
	}

	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		l.logger.Warn("DockerLib Reset ContainerList cause ERROR:", "err", err)
		return false
	}

	var idList []string
	for _, c := range containerList {
		for _, name := range c.Names {
			if strings.HasPrefix(name, prefix) {
				idList = append(idList, c.ID)
				break
			}
		}
	}
	for _, id := range idList {
		l.killByID(ctx, cli, id)
	}
	return true
}

// GetDockerIP 通過容器的名字獲取容器 IP 地址
func (l *DockerLib) GetDockerIP(containerName string) string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		l.logger.Warn("DockerLib Reset NewEnvClient cause ERROR:", "err", err)
	}

	containerID := l.getContainerIDByName(ctx, cli, containerName)
	if containerID == "" {
		l.logger.Warn("DockerLib GetDockerIP no such container ERROR:", "name", containerName)
		return ""
	}
	resp, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		l.logger.Warn("DockerLib GetDockerIP ContainerInspect cause ERROR:", "name", containerName, "err", err)
		return ""
	}

	return resp.NetworkSettings.IPAddress
}

func mapIP(s string) string {
	return strings.Map(func(r rune) rune {
		if r != 46 && (r < 48 || r > 57) { // 只留 [.0-9]
			return -1
		}
		return r
	}, s)
}
