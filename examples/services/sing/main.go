package main

import (
	"github.com/asim/go-micro/v3"
	"sxx-go-micro/common/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/examples/services/sing/handler"
	"sxx-go-micro/plugins/service"
)

// 每秒钟QPS
const QPS = 1

func main() {
	service.Create(
		config.SERVICE_SING,
		func(service micro.Service) {
			// 注册处理函数
			user.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		})
}
