package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/broker/redis/v3"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/asyncRedis/event/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_ = service.Create(
		service.Name(config.SERVICE_ASYNC_EVENT),
		service.CallFunc(func(service micro.Service) {

			pbsb := service.Options().Broker
			err := pbsb.Connect()
			if err != nil {
				fmt.Println("broker connection failed!")
				return
			}

			// 回调
			//func(service micro.Service, pbsb broker.Broker) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service, Pbsb: pbsb},
			)
			//}(service, pbsb)

		}),
		service.Init([]micro.Option{
			micro.Broker(
				// 设置 rocketmq 作为 broker 驱动
				redis.NewBroker(),
			),
		}),
	)
}
