package tcp

import (
	"context"
	"github.com/NetLops/go-imitate-redis/interface/tcp"
	"github.com/NetLops/go-imitate-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

// 感知信号
func ListenAndServeWithSignal(
	cfg *Config,
	handler tcp.Handler) error {

	sigChannel := make(chan os.Signal)
	closeChan := make(chan struct{})
	// 信号通知
	signal.Notify(sigChannel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sigChan := <-sigChannel
		switch sigChan {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	go func() {
		logger.Info("start listen")
	}()
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {
	go func() {
		<-closeChan
		logger.Info("shutting down")
		// 关闭listener、handler
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		// 关闭listener、handler
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go func() {
			logger.Info("accepted link")
		}()
		waitDone.Add(1)
		go func() {
			defer func() { waitDone.Done() }() // 防止 panic
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
