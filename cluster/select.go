package cluster

import "github.com/NetLops/go-imitate-redis/interface/resp"

// execSelect 本地保存即可，发送的时候 会携带dbNum
func execSelect(cluster *ClusterDatabase, con resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(con, cmdArgs)
}
