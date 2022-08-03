package cluster

import "github.com/NetLops/go-imitate-redis/interface/resp"

// 本地执行

func ping(cluster *ClusterDatabase, con resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(con, cmdArgs)
}
