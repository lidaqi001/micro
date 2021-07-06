package Event

import (
	"context"
	"errors"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/server"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/broker/rabbitmq"
	"github.com/streadway/amqp"
)

type New struct {
	// Use for Publish
	Client client.Client

	// Use for Subscribe
	Server server.Server
}

func init() {
	// 设置rabbitmq地址
	rabbitmq.DefaultRabbitURL = helper.GetRabbitmqAddress()
}

// 事件订阅
func (e *New) Subscribe(topic string, queue string, handler interface{}) error {

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
func (e *New) Publish(topic string, ctx context.Context, msg interface{}) error {

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
