package cluster

import (
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/resp/reply"
)

// Del del k1 k2 k3 k4 k5
func Del(cluster *ClusterDatabase, con resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(con, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeArgNumErrReply("IntReply error")
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeErrReply("error: " + errReply.Error())
}
