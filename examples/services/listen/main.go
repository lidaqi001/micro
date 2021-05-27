package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/listen/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_LISTEN,
		func(service micro.Service) {
			// 注册处理函数
			user.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		})
}
