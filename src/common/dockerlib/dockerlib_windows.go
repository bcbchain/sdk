package dockerlib

import (
	"sync"

	"github.com/tendermint/tmlibs/log"
)

// TODO: 爲了照顧 windows 下的編譯，特產生此文件，目前還是空的 api，看怎麼實現方便

// DockerLib 是我們自定義的 Docker API 的 Wrapper
type DockerLib struct {
	logger log.Logger
}

var (
	myLib        *DockerLib
	instanceOnce sync.Once
	initOnce     sync.Once
)

const dockerHubIP = "172.17.0.1"

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
	return dockerHubIP
}

// GetDockerHubIP 獲得本機 Docker 的網卡地址，如果有服務需要 Docker 容器內部訪問，就可以訪問這個地址
func (l *DockerLib) GetDockerHubIP() string {
	return dockerHubIP
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

// Run 運行 Docker 容器，執行某個功能
func (l *DockerLib) Run(dockerImageName, containerName string, params *DockerRunParams) bool {
	l.logger.Info("DockerLib Run", "image", dockerImageName, "containerName", containerName, "params", params)
	return true
}

// Kill 殺死一個 Docker 容器，並且清理現場
func (l *DockerLib) Kill(containerName string) bool {
	return true
}

// Status 查詢一個容器的狀態
func (l *DockerLib) Status(containerName string) bool {
	return true
}

// Reset 殺掉所有自己啓動的容器(以特定字冠命名的)
func (l *DockerLib) Reset(prefix string) bool {
	return true
}

// GetDockerIP 通過容器的名字獲取容器 IP 地址
func (l *DockerLib) GetDockerIP(containerName string) string {
	return dockerHubIP
}
