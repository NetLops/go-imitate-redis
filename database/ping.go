package database

import (
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
