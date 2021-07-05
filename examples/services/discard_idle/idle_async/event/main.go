package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/idle_async/event/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_ = service.Create(
		service.Name(config.SERVICE_ASYNC_EVENT),
		service.CallFunc(func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service},
			)
		}),
	)
}
