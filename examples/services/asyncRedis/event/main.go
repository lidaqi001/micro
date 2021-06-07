package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/broker/redis/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/asyncRedis/event/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	service.Create(config.SERVICE_ASYNC_EVENT, func(service micro.Service) {

		pbsb := service.Options().Broker
		err := pbsb.Connect()
		if err != nil {
			fmt.Println("broker connection failed!")
			return
		}
		fmt.Println(err)
		fmt.Println(service.Options().Broker.Address())
		//defer pbsb.Disconnect()

		// 回调
		func(service micro.Service, pbsb broker.Broker) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service, Pbsb: pbsb},
			)
		}(service, pbsb)

	}, micro.Broker(
		// 设置 rocketmq 作为 broker 驱动
		redis.NewBroker()),
	)
}
