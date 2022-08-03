package tcp

import (
	"bufio"
	"context"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"github.com/NetLops/go-imitate-redis/lib/sync/atomic"
	"github.com/NetLops/go-imitate-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean // 都要支持并发
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// 关闭客户端
func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

/* 业务 */
func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	// 所有连接的客户
	handler.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for true {
		msg, err := reader.ReadString('\n') // '\n'为标志位
		if err != nil {
			// 结束符 // 代表客户端退出
			if err == io.EOF {
				logger.Info("Connectiong close")
				handler.activeConn.Delete(client)
			} else {
				go func() {
					logger.Warn(err)
				}()
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}

}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)

	// 关闭所有客户端
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}
