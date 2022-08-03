package cluster

import (
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/resp/reply"
)

// flushdb 广播群发
func flushdb(cluster *ClusterDatabase, con resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(con, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("error: " + errReply.Error())
}
