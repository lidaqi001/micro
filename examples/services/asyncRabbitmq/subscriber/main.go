package main

import (
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/server"
	"github.com/lidaqi001/micro/common/config"
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
			// 问题：因为每一次队列名称都是随机生成的，所以会有下面的问题
			// 场景：如果当前订阅者异常退出，下一次再注册订阅者时，中间这部分消息会丢失

			//b1 := "singEvent"
			//_ = micro.RegisterSubscriber(
			//	b1,
			//	service.Server(),
			//	handler.CallSing,
			//)

			b1 := config.EVENT_A
			brkrSub := broker.NewSubscribeOptions(
				broker.Queue(config.QUEUE_A),
				broker.DisableAutoAck(),
				rabbitmq.DurableQueue(),
			)
			_ = micro.RegisterSubscriber(
				b1,
				service.Server(),
				handler.SingEvent,
				server.SubscriberContext(brkrSub.Context),
				server.SubscriberQueue(brkrSub.Queue),
			)

			// 指定队列
			// 队列名称不能重复！！！
			// 将队列名称放在一个配置文件中，以队列名称定义常量名

			b2 := config.EVENT_A
			brkrSub2 := broker.NewSubscribeOptions(
				broker.Queue(config.QUEUE_B),
				broker.DisableAutoAck(),
				rabbitmq.DurableQueue(),
			)
			_ = micro.RegisterSubscriber(
				b2,
				service.Server(),
				handler.SingEvent,
				server.SubscriberContext(brkrSub2.Context),
				server.SubscriberQueue(brkrSub2.Queue),
			)

			b3 := config.EVENT_B
			brkrSub3 := broker.NewSubscribeOptions(
				broker.Queue(config.QUEUE_C),
				broker.DisableAutoAck(),
				rabbitmq.DurableQueue(),
			)
			// 注册订阅
			_ = micro.RegisterSubscriber(
				b3,
				service.Server(),
				handler.CallSing,
				server.SubscriberContext(brkrSub3.Context),
				server.SubscriberQueue(brkrSub3.Queue),
			)
		},
	)
}
