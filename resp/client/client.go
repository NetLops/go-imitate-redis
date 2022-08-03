package client

import (
	"fmt"
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"github.com/NetLops/go-imitate-redis/lib/sync/wait"
	"github.com/NetLops/go-imitate-redis/resp/parser"
	"github.com/NetLops/go-imitate-redis/resp/reply"
	"net"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
)

type Client struct {
	conn        net.Conn      // 与服务器的tcp的连接
	pendingReqs chan *Request // 等待发送的请求
	waitingReqs chan *Request // 等待服务器响应的请求
	ticker      *time.Ticker  // 用于触发心跳包的计时器
	addr        string

	//ctx        context.Context
	//cancelFunc context.CancelFunc
	//status  int32
	working *sync.WaitGroup // 有请求正在处理不能立即停止，用于实现 graceful shutdown
}

type Request struct {
	id        uint64     // 	请求id
	args      [][]byte   // 上行参数
	reply     resp.Reply // 收到的返回值
	heartbeat bool       // 标记是否是心跳请求
	waiting   *wait.Wait // 调用协程发送请求后通过 waitGroup 等待请求异步处理完成
	err       error
}

const (
	chanSize = 256
	maxWait  = 3 * time.Second
)

// MakeClient creates a new client
func MakeClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr:        addr,
		conn:        conn,
		pendingReqs: make(chan *Request, chanSize),
		waitingReqs: make(chan *Request, chanSize),
		working:     &sync.WaitGroup{},
	}, nil
}

// Start starts asynchronous goroutines
func (client *Client) Start() {
	client.ticker = time.NewTicker(10 * time.Second)
	go client.handleWrite()
	go func() {
		err := client.handleRead()
		if err != nil {
			logger.Error(err)
		}
	}()
	go client.heartbeat()
}

// Close stops asynchronous goroutines and close connection
func (client Client) Close() {
	client.ticker.Stop()
	// stop new request
	close(client.pendingReqs) // 先阻止新请求进入队列

	// wait stop process
	client.working.Wait()

	// clean
	_ = client.conn.Close()   // 关闭与服务端的连接，连接关闭后读协程会退出
	close(client.waitingReqs) // 关闭队列

}

func (client *Client) handleConnectionError(err error) error {
	err1 := client.conn.Close()
	if err1 != nil {
		if opError, ok := err1.(*net.OpError); ok {
			if opError.Err.Error() != "use of closed network connection" { // 使用了已经关闭的网络连接
				return err1
			}
		} else {
			return err1
		}
	}
	conn, err1 := net.Dial("tcp", client.addr)
	if err1 != nil {
		logger.Error(err1)
		return err1
	}
	client.conn = conn
	go func() {
		_ = client.handleRead()
	}()
	return nil
}

func (client *Client) heartbeat() {
	for range client.ticker.C {
		client.doHeartbeat()
	}
}

func (client *Client) doHeartbeat() {
	request := &Request{
		args:      [][]byte{[]byte("PING")},
		heartbeat: true,
		waiting:   &wait.Wait{},
	}
	request.waiting.Add(1)
	client.working.Add(1)
	defer client.working.Done()
	client.pendingReqs <- request
	request.waiting.WaitWithTimeout(maxWait)
}

// Send
// 调用者将请求发送给后台协程，并通过 wait group 等待异步处理完成
func (client *Client) Send(args [][]byte) resp.Reply {
	request := &Request{
		args:      args,
		heartbeat: false,
		waiting:   &wait.Wait{},
	}
	request.waiting.Add(1)
	client.working.Add(1)
	defer client.working.Done()
	client.pendingReqs <- request                       // 	请求入队
	timeout := request.waiting.WaitWithTimeout(maxWait) // 等待响应或者超时
	if timeout {
		return reply.MakeErrReply("server time out")
	}
	if request.err != nil {
		return reply.MakeErrReply("request failed")
	}
	return request.reply // 服务端的响应
}

// handleWrite
// 后台的读写协程
func (client *Client) handleWrite() {
	for req := range client.pendingReqs {
		client.doRequest(req)
	}
}

// doRequest 发送请求
func (client *Client) doRequest(req *Request) {
	if req == nil || len(req.args) == 0 {
		return
	}
	// 序列化请求
	re := reply.MakeMultiBulkReply(req.args)
	//fmt.Println(string(re.ToBytes()))
	bytes := re.ToBytes()
	_, err := client.conn.Write(bytes)
	i := 0
	// 失败重试机制（重试次数3）
	for err != nil && i < 3 {
		err := client.handleConnectionError(err)
		if err == nil {
			_, err = client.conn.Write(bytes)
		}
		i++
	}
	if err == nil {
		// 发送成功等待服务器响应
		client.waitingReqs <- req
	} else {
		req.err = err
		req.waiting.Done()
	}
}

func (client *Client) finishRequest(reply resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			logger.Error(err)
		}
	}()
	request := <-client.waitingReqs
	if request == nil {
		return
	}

	request.reply = reply

	//fmt.Println(request.reply)
	if request.waiting != nil {
		request.waiting.Done()
	}
}

func (client *Client) handleRead() error {

	ch := parser.ParseStream(client.conn)
	for payload := range ch {
		if payload.Err != nil {
			client.finishRequest(reply.MakeErrReply(payload.Err.Error()))
			continue
		}
		data := payload.Data
		fmt.Println(data, reflect.TypeOf(data))
		client.finishRequest(data)
	}
	return nil
}
