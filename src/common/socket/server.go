package socket

import (
	"blockchain/types"
	"bufio"
	"container/list"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tendermint/tmlibs/log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type CallBackFunc func(map[string]interface{}) (interface{}, error)

// Server server information about socket
type Server struct {
	timeout    time.Duration
	listenAddr string
	listener   net.Listener
	methods    map[string]CallBackFunc

	connList *list.List
	mtx      sync.Mutex
	logger   log.Logger
}

// NewServer newServer to create server object about socket and listen client connection
func NewServer(listenAddr string, methods map[string]CallBackFunc, timeout time.Duration, logger log.Logger) (svr *Server, err error) {

	logger.Info(fmt.Sprintf("New server with listenaddr=%s, methods=%v, timeout=%d", listenAddr, methods, timeout))
	server := Server{
		listenAddr: listenAddr,
		methods:    methods,
		timeout:    timeout,
		logger:     logger,
		connList:   list.New()}

	err = server.listen()
	if err != nil {
		return
	}

	return &server, nil
}

// Start start a routine to accept new connection and create routine to operate it
func (svr *Server) Start() (err error) {

	// notify system signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for sig := range c {
			next := svr.connList.Front()
			for next != nil {
				next.Value.(net.Conn).Close()
				next = next.Next()
			}
			svr.logger.Info("captured %v, exiting...\n", sig)
			os.Exit(1)
		}
	}()

	for {
		cliConn, err := svr.listener.Accept()
		if err != nil {
			return err
		}

		// save connection
		svr.mtx.Lock()
		svr.connList.PushBack(cliConn)
		svr.mtx.Unlock()

		svr.logger.Info("Accept new connection", "RemoteAddr", cliConn.RemoteAddr())
		go svr.readRequest(cliConn)
	}
}

func (svr *Server) listen() (err error) {
	proto, address := "tcp", svr.listenAddr
	parts := strings.SplitN(svr.listenAddr, "://", 2)
	if len(parts) == 2 {
		proto = parts[0]
		address = parts[1]
	}

	svr.listener, err = net.Listen(proto, address)

	return
}

func (svr *Server) readRequest(conn net.Conn) {
	defer svr.close(conn)

	for {
		value, err := readMessage(conn)
		if err != nil {
			return
		}
		var req = &Request{}
		err = jsoniter.Unmarshal(value, req)
		if err != nil {
			return
		}
		//svr.logger.Info("NewRequest", "value", fmt.Sprintf("%v", req))

		go svr.handleRequest(req, conn)
	}
}

func (svr *Server) handleRequest(req *Request, conn net.Conn) {
	defer serverRecover(conn, req)

	method := svr.methods[req.Method]
	if method == nil {
		panic("Invalid method")
	}

	res, err := method(req.Data)
	if err != nil {
		panic(err)
	}

	svr.logger.Debug(fmt.Sprintf("handlerRequest req=%v result", req), "res", res)
	var resp Response
	resp.Code = types.CodeOK
	resp.Log = "ok"
	resp.Result.Index = req.Index
	resp.Result.Method = req.Method
	resp.Result.Data = res

	svr.mtx.Lock()
	defer svr.mtx.Unlock()
	w := bufio.NewWriter(conn)
	err = writeMessage(resp, w)
	if err != nil {
		panic(err)
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}
}

func (svr *Server) close(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		svr.logger.Info("Close connection error: " + err.Error())
	}

	svr.mtx.Lock()
	defer svr.mtx.Unlock()
	next := svr.connList.Front()
	for next != nil {
		if next.Value.(net.Conn) == conn {
			svr.connList.Remove(next)
			break
		}
		next = next.Next()
	}
}
