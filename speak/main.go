package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/Common/service"
	proto "sxx-go-micro/proto"
	"sxx-go-micro/speak/handler"
)

func main() {
	service.CreateService(
		context.Background(),
		config.SERVICE_SPEAK,
		func(service micro.Service) {
			// 注册处理函数
			proto.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		})
}
