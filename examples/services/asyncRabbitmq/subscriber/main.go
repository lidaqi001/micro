package main

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/server"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/asyncRabbitmq/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/rabbitmqPack"
)

func main() {

	rabbitmqPack.Create(

		config.SERVICE_ASYNC_SUBSCRIBER,

		func(service micro.Service, pbsb broker.Broker) {
			// 注册订阅

			// 不指定队列
			// 每次运行订阅服务，默认会分配一个随机队列

			b1 := "singEvent"
			_ = micro.RegisterSubscriber(
				b1,
				service.Server(),
				handler.CallSing,
			)

			// 指定队列
			// 队列名称不能重复！！！
			// 将队列名称放在一个配置文件中，以队列名称定义常量名

			b2 := "singEvent"
			brkrSub := broker.NewSubscribeOptions(
				broker.Queue("t2"),
				broker.DisableAutoAck(),
				rabbitmq.DurableQueue(),
			)
			_ = micro.RegisterSubscriber(
				b2,
				service.Server(),
				handler.SingEvent,
				server.SubscriberContext(brkrSub.Context),
				server.SubscriberQueue(brkrSub.Queue),
			)

			b3 := "callSing"
			brkrSub3 := broker.NewSubscribeOptions(
				broker.Queue("t3"),
				broker.DisableAutoAck(),
				rabbitmq.DurableQueue(),
			)
			// 注册订阅
			_ = micro.RegisterSubscriber(
				b3,
				service.Server(),
				handler.SingEvent,
				server.SubscriberContext(brkrSub3.Context),
				server.SubscriberQueue(brkrSub3.Queue),
			)
		},
	)
}
