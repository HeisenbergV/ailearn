package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fogate/internal/config"
	"fogate/internal/handler"
	"fogate/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/fogate-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 启动 HTTP 服务
	go func() {
		fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
		server.Start()
	}()

	// 启动终端服务
	go func() {
		ctx.Terminal.Start()
	}()

	// 等待中断信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
