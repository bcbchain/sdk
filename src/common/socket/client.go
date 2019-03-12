package socket

import (
	"blockchain/types"
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"math"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Client client information about socket
type Client struct {
	reqSent *list.List

	addr             string
	conn             net.Conn
	timeout          time.Duration
	disableKeepAlive bool

	counter uint64
	mtx     sync.Mutex
	logger  log.Logger
}

// NewClient newClient to create socket client object and connect to server
func NewClient(addr string, timeout time.Duration, disableKeepAlive bool, logger log.Logger) (cli *Client, err error) {
	if timeout == 0 {
		timeout = 60
	}

	cli = &Client{
		reqSent:          list.New(),
		addr:             addr,
		timeout:          timeout,
		disableKeepAlive: disableKeepAlive,
		counter:          0,
		logger:           logger}

	err = cli.connect()
	if err != nil {
		return
	}

	go cli.recvResponseRoutine()

	return
}

// Call call service with method and data
func (cli *Client) Call(method string, data map[string]interface{}) (value interface{}, err error) {

	cli.logger.Debug("Calling...", "data", data)
	req := Request{Method: method, Data: data, Index: cli.index()}
	if cli.disableKeepAlive {
		defer cli.conn.Close()
	}
	cli.logger.Debug("Calling2...", "index", req.Index)

	// send request
	w := bufio.NewWriter(cli.conn)
	cli.mtx.Lock()
	err = writeMessage(req, w)
	if err != nil {
		cli.mtx.Unlock()
		cli.logger.Error(err.Error())
		return
	}
	err = w.Flush()
	if err != nil {
		cli.mtx.Unlock()
		return
	}
	cli.mtx.Unlock()

	// wait response
	respChan := make(chan *Response, 1)
	cli.sentReq(req.Index, respChan)

	// notify system signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("captured %v, exiting...\n", sig)
			os.Exit(1)
		}
	}()

	cli.logger.Debug(fmt.Sprintf("Call select, %d", cli.timeout))
	select {
	case sig := <-c:
		return nil, errors.New(fmt.Sprintf("captured %v", sig))
	case <-time.After(cli.timeout * time.Second):
		return nil, errors.New("Recv time out ")
	case resp := <-respChan:
		//resp := <-respChan
		if resp.Code == types.CodeOK {
			return resp.Result.Data, nil
		} else {
			return nil, errors.New(resp.Log)
		}
	}
}

// Close close connection
func (cli *Client) Close() (err error) {
	err = cli.conn.Close()
	if err != nil {
		return
	}
	cli.conn = nil

	return
}

func (cli *Client) connect() (err error) {
	proto, address := "tcp", cli.addr
	parts := strings.SplitN(cli.addr, "://", 2)
	if len(parts) == 2 {
		proto, address = parts[0], parts[1]
	}

	var keepAlive time.Duration
	if cli.disableKeepAlive == false {
		keepAlive = 5 * time.Second
	}

	dialer := net.Dialer{Timeout: cli.timeout * time.Second, KeepAlive: keepAlive}
	cli.conn, err = dialer.Dial(proto, address)
	if err != nil {
		return err
	}

	return
}

func (cli *Client) recvResponseRoutine() {
	for {
		value, err := readMessage(cli.conn)
		if err != nil {
			return
		}

		var resp Response
		err = jsoniter.Unmarshal(value, &resp)
		if err != nil {
			return
		}

		go cli.didRecvResponse(resp)
	}
}

func (cli *Client) didRecvResponse(resp Response) {
	tryCount := 3
	var next *list.Element
	for tryCount > 0 {
		next = cli.reqSent.Front()
		for next != nil {
			if next.Value.(ReqResp).Index == resp.Result.Index {
				break
			}

			next = next.Next()
		}

		if next != nil {
			break
		}
		tryCount--
		time.Sleep(100 * time.Microsecond)
	}

	if next != nil {
		next.Value.(ReqResp).RespChan <- &resp
		cli.reqSent.Remove(next)
	} else {
		cli.logger.Error("didRecvResponse", "response index", resp.Result.Index, "reqSent", cli.reqSent)
	}
}

func (cli *Client) sentReq(index uint64, respChan chan *Response) {
	cli.mtx.Lock()
	defer cli.mtx.Unlock()
	cli.reqSent.PushBack(ReqResp{Index: index, RespChan: respChan})
}

func (cli *Client) index() uint64 {
	cli.mtx.Lock()
	defer cli.mtx.Unlock()
	cli.counter += 1
	return cli.counter % math.MaxUint64
}
