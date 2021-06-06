package main

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/async/event/handler"
	"github.com/lidaqi001/micro/plugins/service"
	"os"
)

func main() {
	rabbitmq.DefaultRabbitURL = "amqp://sxx:sxx123456@192.168.1.146:5672"
	//rabbitmq.DefaultRabbitURL = getAddr()
	service.Create(
		config.SERVICE_ASYNC_EVENT,
		func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(
				service.Server(), &handler.DemoServiceHandler{Service: service},
			)
		}, micro.Broker(
			rabbitmq.NewBroker(rabbitmq.DurableExchange()),
		))
}

func getAddr() string {
	e := os.Getenv("BROKER_ADDR")
	if len(e) == 0 {
		return ""
	}
	return e
}
