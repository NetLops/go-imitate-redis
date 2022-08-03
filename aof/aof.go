package aof

import (
	"github.com/NetLops/go-imitate-redis/config"
	databaseinterface "github.com/NetLops/go-imitate-redis/interface/database"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"github.com/NetLops/go-imitate-redis/lib/utils"
	"github.com/NetLops/go-imitate-redis/resp/connection"
	"github.com/NetLops/go-imitate-redis/resp/parser"
	"github.com/NetLops/go-imitate-redis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine = [][]byte

const (
	aofQueueSize = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// AofHandler receive msgs from channel and write to AOF file
type AofHandler struct {
	db databaseinterface.Database
	//database    database.StandaloneDatabase
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

// NewAOFHandler creates a new aof.AofHandler
func NewAofHandler(database databaseinterface.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.db = database
	// LoadAof
	handler.LoadAof()

	aofFile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	// Channel
	handler.aofChan = make(chan *payload, aofQueueSize)
	go func() {
		handler.handleAof()
	}()

	return handler, nil
}

// AddAof send command to aof goroutine through channel
func (handler *AofHandler) AddAof(dbIndex int, cmdLine CmdLine) {
	// 判断是否开启持久化功能
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}

}

// handleAof listen aof channel and write into file
// payload(set k v) <- aofChan
func (handler *AofHandler) handleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}
	}

}

// LoadAof read aof file
func (handler *AofHandler) LoadAof() {
	open, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		err := open.Close()
		if err != nil {
			return
		}
	}()
	ch := parser.ParseStream(open)
	fackConn := &connection.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF { // 文件结束符 就退出
				break
			}
			logger.Error(p.Err)
			//TODO 出现错误继续
			continue
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		result, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi mulk")
			continue
		}
		resp := handler.db.Exec(fackConn, result.Args)
		if reply.IsErrReply(resp) {
			logger.Error("error err", resp.ToBytes())
		}
	}
}
