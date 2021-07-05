package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/services/current_async_rabbitmq/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/event"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {

	_ = service.Create(
		service.Name(config.SERVICE_ASYNC_SUBSCRIBER),
		service.RabbitmqBroker(true),
		service.CallFunc(func(service micro.Service) {

			// 注册订阅

			/**
			指定队列（使用，稳定的线上运行时）

			注意点：队列名称不能重复！！！

			将队列名称放在一个配置文件中，以队列名称定义常量名
			确保队列名称没有重复
			*/

			e := Event.New{Server: service.Server()}

			_ = e.Subscribe(config.EVENT_A, config.QUEUE_A, handler.SingEvent)

			_ = e.Subscribe(config.EVENT_A, config.QUEUE_B, handler.SingEvent)

			_ = e.Subscribe(config.EVENT_B, config.QUEUE_C, handler.CallSing)

			/**
			不指定队列（不使用，会导致消息数据丢失）

			每次运行订阅服务，默认会分配一个随机队列
			问题：因为每一次队列名称都是随机生成的，所以会有下面的问题
			场景：如果当前订阅者异常退出，下一次再注册订阅者时，中间这部分消息会丢失
			*/

			/*
				b1 := "singEvent"
				_ = micro.RegisterSubscriber(
					b1,
					service.Server(),
					handler.CallSing,
				)
			*/
		}),
	)
}
