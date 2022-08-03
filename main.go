package main

import (
	"fmt"
	"github.com/NetLops/go-imitate-redis/config"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"github.com/NetLops/go-imitate-redis/resp/handler"
	"github.com/NetLops/go-imitate-redis/tcp"
	"os"
)

const configFile = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, handler.MakeHandler())
	if err == nil {
		logger.Error(err)
	}
}

//
//func main() {
//	b := []byte("test\r\n")
//	fmt.Println(string(b))
//
//
//}
//"*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
//"*2\r\n$3\r\nGET\r\n$2\r\nt3\r\n"
//"*2\r\n$3\r\nGET\r\n$2\r\nt4\r\n"
//"*2\r\n$6\r\nSELECT\r\n$1\r\n3\r\n"
//"*2\r\n$4\r\nkeys\r\n$6\r\n[a-z]*\r\n"
//
//"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
//"*3\r\n$3\r\nSET\r\n$2\r\nt3\r\n$2\r\n34\r\n"
//"*3\r\n$3\r\nSET\r\n$2\r\nt4\r\n$2\r\n14\r\n"
