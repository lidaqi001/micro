package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/broker/kafka/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/(discard)asyncKafka/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_ = service.Create(
		service.Name(config.SERVICE_ASYNC_SUBSCRIBER),
		service.CallFunc(func(service micro.Service) {

			pbsb := service.Options().Broker

			if err := pbsb.Connect(); err != nil {
				fmt.Println("broker connection failed!")
				return
			}
			//defer pbsb.Disconnect()

			// 回调
			func(service micro.Service, pbsb broker.Broker) {

				_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.CallSing)

				_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.SingEvent)
			}(service, pbsb)

		}),
		service.Init([]micro.Option{
			micro.Broker(
				kafka.NewBroker(),
			),
		}),
	)
}
