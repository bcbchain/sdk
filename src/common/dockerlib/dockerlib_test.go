package dockerlib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tmlibs/log"
)

func TestDockerLib_GetDockerHubIP(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	ip := lib.GetDockerHubIP()
	println(ip)
}

func TestDockerLib_Run(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	lib.Kill("my8080")
	params := DockerRunParams{
		PortMap: map[string]HostPort{
			"8000": {
				Port: "8080",
				Host: "0.0.0.0",
			},
		},
		WorkDir:    "/",
		AutoRemove: true,
		Cmd:        []string{"sh", "-c", "python2 -m SimpleHTTPServer"},
	}
	ret := lib.Run("python:2-alpine", "my8080", &params)
	assert.Equal(t, ret, true)

	timer := time.NewTimer(10 * time.Second)
	checkTimer := time.NewTicker(20 * time.Millisecond)
	defer func() { checkTimer.Stop() }()
	count := 0
	for {
		select {
		case <-checkTimer.C:
			resp, err := http.Get("http://localhost:8080/")
			if err == nil && resp.StatusCode == 200 {
				goto GOTOEND
			}
			fmt.Println("err =", err, "; resp =", resp)
			count++
			continue
		case <-timer.C:
			assert.Error(t, fmt.Errorf("啓動時間過長"), "")
			goto GOTOEND
		}
	}
GOTOEND:
	fmt.Println("count=", count)
	resp, err := http.Get("http://localhost:8080/")
	assert.Equal(t, err, nil)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}

func TestDockerLib_Run2(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	lib.Kill("my8000")
	params := DockerRunParams{
		WorkDir: "/",
		Cmd:     []string{"sh", "-c", "python2 -m SimpleHTTPServer"},
	}
	ret := lib.Run("python:2-alpine", "my8000", &params)
	assert.Equal(t, ret, true)

	ip := lib.GetDockerIP("my8000")
	fmt.Println("ip=", ip)
	timer := time.NewTimer(10 * time.Second)
	checkTimer := time.NewTicker(100 * time.Millisecond)
	defer func() { checkTimer.Stop() }()
	count := 0
	for {
		select {
		case <-checkTimer.C:
			resp, err := http.Get("http://" + ip + ":8000/")
			if err == nil && resp.StatusCode == 200 {
				goto GOTOEND
			}
			fmt.Println("err =", err, "; statusCode =", resp.StatusCode)
			count++
			continue
		case <-timer.C:
			assert.Error(t, fmt.Errorf("啓動時間過長"), "")
			goto GOTOEND
		}
	}
GOTOEND:
	fmt.Println("count=", count)
	resp, err := http.Get("http://" + ip + ":8000/")
	assert.Equal(t, err, nil)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
func TestDockerLib_GetDockerIP(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	ip := lib.GetDockerIP("my8000")
	fmt.Println(ip)
	assert.Equal(t, ip, "172.17.0.3")
}

func TestDockerLib_Kill(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	result := lib.Kill("my8000")
	assert.Equal(t, result, true)
}

func TestDockerLib_GetMyIntranetIP(t *testing.T) {
	logger := log.NewOldTMLogger(os.Stdout)
	lib := GetDockerLib()
	lib.Init(logger)
	ip := lib.GetMyIntranetIP()
	assert.Equal(t, ip, "192.168.1.224")
}
