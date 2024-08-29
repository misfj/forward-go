package main

import (
	"context"
	"forward-go/config"
	"forward-go/core"
	"forward-go/global"
	"forward-go/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//初始化配置文件,加载配置文件到内存。如果不存在配置文件,生成默认的配置文件。
	config.LoadConfig()
	//初始化各种组件连接
	global.LoadGlobal()
	server := core.New()
	err := server.Init()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	newCtx, cancel := context.WithCancel(context.Background())
	server.Run(newCtx)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	sig := <-interrupt
	log.Info("recv signal:", sig)

	cancel()
	server.Close()
	server.Wait()
}
