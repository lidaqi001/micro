package rabbitmqPack

import (
	"context"
	"errors"
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/server"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/broker/rabbitmq"
	"github.com/streadway/amqp"
	//"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/lidaqi001/micro/plugins/service"
)

func Create(serviceName string, registerService func(service micro.Service, pbsb broker.Broker)) {

	// 设置rabbitmq地址
	rabbitmq.DefaultRabbitURL = helper.GetConfig("RABBITMQ_ADDR", config.RABBITMQ_ADDR)

	_ = service.Create(
		service.Name(serviceName),
		service.CallFunc(func(service micro.Service) {

			pbsb := service.Options().Broker

			if err := pbsb.Connect(); err != nil {
				fmt.Println("broker connection failed!")
				logger.Fatal("broker connection failed!")
				//return
			}

			// 回调
			registerService(service, pbsb)

		}),
		service.Init([]micro.Option{
			micro.Broker(
				// 设置 rabbitmq 为 broker 驱动
				rabbitmq.NewBroker(
					// 设置：Exchange 为持久化
					// If this option is not set, the exchange will be deleted when rabbitmq restarts
					rabbitmq.DurableExchange(),
					// 设置：订阅时创建持久化队列
					//rabbitmq.PrefetchGlobal(),
				),
			),
		}),
	)
}

type Event struct {
	// Use for Publish
	Client client.Client

	// Use for Subscribe
	Server server.Server
}

// 事件订阅
func (e *Event) Subscribe(topic string, queue string, handler interface{}) error {

	if e.Server == nil {
		return errors.New("Event.Server cannot be nil!")
	}

	options := broker.NewSubscribeOptions(
		// 指定队列
		broker.Queue(queue),
		// 禁止自动Ack
		broker.DisableAutoAck(),
		// 持久化队列
		rabbitmq.DurableQueue(),
	)

	return micro.RegisterSubscriber(
		topic,
		e.Server,
		handler,
		server.SubscriberContext(options.Context),
		server.SubscriberQueue(options.Queue),
	)
}

//事件发布
func (e *Event) Publish(topic string, ctx context.Context, msg interface{}) error {

	if e.Client == nil {
		return errors.New("Event.Client cannot be nil!")
	}

	// 设置DeliveryMode 持久化消息
	publishOptsCtx := context.WithValue(context.Background(), rabbitmq.DeliveryMode{}, amqp.Persistent)

	// 设置Priority 消息优先级
	//publishOptsCtx = context.WithValue(ctx, rabbitmq.Priority{}, 0)

	p := micro.NewEvent(topic, e.Client)
	if err := p.Publish(ctx, msg, client.PublishContext(publishOptsCtx)); err != nil {
		logger.Error("[pub] failed: %v", err)
		return err
	}

	return nil
}
