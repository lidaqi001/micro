package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/asyncRocketmq/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/asyncRocketmq"
)

func main() {
	asyncRocketmq.Create(
		config.SERVICE_ASYNC_SUBSCRIBER_ROCKETMQ,
		func(service micro.Service, pbsb broker.Broker) {
			pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.SingEvent)
			//_ = micro.RegisterSubscriber("singEvent", service.Server(), handler.SingEvent)
			//_ = micro.RegisterSubscriber("callSing", service.Server(), handler.CallSing)

		})
}
