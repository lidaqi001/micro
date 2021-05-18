package main

import (
	"github.com/asim/go-micro/v3"
	"sxx-go-micro/examples/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/examples/services/listen/handler"
	"sxx-go-micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_LISTEN,
		func(service micro.Service) {
			// 注册处理函数
			user.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		})
}
