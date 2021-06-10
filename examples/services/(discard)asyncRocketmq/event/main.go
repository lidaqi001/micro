package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/(discard)asyncRocketmq/event/handler"
	"github.com/lidaqi001/micro/plugins/rocketmqPack"
)

func main() {
	rocketmqPack.Create(
		config.SERVICE_ASYNC_EVENT,
		func(service micro.Service, pbsb broker.Broker) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service, Pbsb: pbsb},
			)
		})
}
