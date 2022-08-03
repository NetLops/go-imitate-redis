package database

import (
	databaseInterface "github.com/NetLops/go-imitate-redis/interface/database"
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connection, args databaseInterface.CmdLine) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {
}
