package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/sing/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_SING,
		func(service micro.Service) {
			logger.Log(logger.DebugLevel, "ceshi")
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(service.Server(), &handler.DemoServiceHandler{})
		})
}
