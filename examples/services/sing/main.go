package main

import (
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/sing/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	err := service.Create(
		service.Name(config.SERVICE_SING),
		service.CallFunc(func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(service.Server(), &handler.DemoServiceHandler{})
		}),
	)
	fmt.Println(err)
}
