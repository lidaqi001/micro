package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/speak/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_ = service.Create(
		service.Name(config.SERVICE_SPEAK),
		service.CallFunc(func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		}),
	)
}
