package cluster

import (
	"github.com/NetLops/go-imitate-redis/interface/resp"
	"github.com/NetLops/go-imitate-redis/resp/reply"
)

// Rename / rename k1 k2
func Rename(cluster *ClusterDatabase, con resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1])  // k1
	dest := string(cmdArgs[2]) // k2

	srcPeer := cluster.peerPicker.PickNode(src) // xxx.xxx.xxx.xxx
	destPeer := cluster.peerPicker.PickNode(dest)
	//TODO 这边逻辑可以改成 删除原key 然后 再将新的key插入
	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, con, cmdArgs)
}
