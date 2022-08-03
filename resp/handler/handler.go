package handler

import (
	"context"
	"github.com/NetLops/go-imitate-redis/cluster"
	"github.com/NetLops/go-imitate-redis/config"
	"github.com/NetLops/go-imitate-redis/database"
	databaseInterface "github.com/NetLops/go-imitate-redis/interface/database"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"github.com/NetLops/go-imitate-redis/lib/sync/atomic"
	"github.com/NetLops/go-imitate-redis/resp/connection"
	"github.com/NetLops/go-imitate-redis/resp/parser"
	"github.com/NetLops/go-imitate-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseInterface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db databaseInterface.Database
	//TODO: 实现Database
	//db = &database.EchoDatabase{}

	if config.Properties.Self != "" && len(config.Properties.Peers) > 0 {
		db = cluster.MakeClusterDatabase() // 集群版
	} else {
		db = database.NewStandaloneDatabase() // 单机版
	}
	//db = database.NewStandaloneDatabase() // 单机模式
	return &RespHandler{
		db: db,
	}
}

// closeClient 关闭一个client
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

// Handle 处理tcp连接
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, 1)
	// 交给Parser
	ch := parser.ParseStream(conn)
	// 监听管道
	for payload := range ch {
		// 如果有错误
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") { // 客户端做了挥手
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// 协议错误	protocol error
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes()) // 回写出错
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		// Exec
		// 用户发送过来的指令为空
		if payload.Data == nil {
			continue
		}
		res, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, res.Args)

		if result != nil {
			_ = client.Write(result.ToBytes())
			//_ = client.Write(reply.MakeOkReply().ToBytes())
			//_ = client.Write(reply.MakeBulkReply([]byte("test")).ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

// Close 关闭整个redis
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	// 标志位为0
	r.closing.Set(true)
	// 关闭所有客户端
	r.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	r.db.Close()
	return nil
}
