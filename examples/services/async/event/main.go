package main

import (
	"github.com/asim/go-micro/v3"
	"sxx-go-micro/common/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/examples/services/async/event/handler"
	"sxx-go-micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_ASYNC_EVENT,
		func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service},
			)
		})
}
