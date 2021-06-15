package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/current_asyncRabbitmq/event/handler"
	"github.com/lidaqi001/micro/plugins/rabbitmqPack"
)

func main() {

	rabbitmqPack.Create(

		config.SERVICE_ASYNC_EVENT,

		func(service micro.Service, pbsb broker.Broker) {

			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service},
			)

		},
	)
}
