package database

import "github.com/NetLops/go-imitate-redis/interface/resp"

type CmdLine = [][]byte

// Database Redis 业务核心
type Database interface {
	Exec(client resp.Connection, args CmdLine) resp.Reply // 执行
	Close()                                               // 关闭
	AfterClientClose(c resp.Connection)                   // 客户端可能需要清理
}

// DataEntity 指代任何的数据结构
type DataEntity struct {
	Data interface{}
}
