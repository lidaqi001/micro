package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/codec/json"
	"os"
	"strconv"
)

type rBroker struct {
	addrs     []string
	opts      broker.Options
	connected bool
}
type publication struct {
}

func init() {
	cmd.DefaultBrokers["rabbitmq"] = NewBroker
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// default to json codec
		Codec:   json.Marshaler{},
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9876"}
	}

	return &rBroker{
		addrs: cAddrs,
		opts:  options,
	}
}

func (r *rBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	var cAddrs []string
	for _, addr := range r.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9876"}
	}
	r.addrs = cAddrs
	return nil
}

func (r *rBroker) Options() broker.Options {
	return r.opts
}

func (r *rBroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return "http://test.mqrest.cn-hangzhou.aliyuncs.com"
}

func (r *rBroker) Connect() error {
	panic("implement me")
}

func (r *rBroker) Disconnect() error {
	panic("implement me")
}

func (r *rBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {

	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(r.addrs)),
		producer.WithRetry(2),
		producer.WithGroupName("GID_xxxxxx"),
		//producer.WithQueueSelector(producer.NewManualQueueSelector())
	)

	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
	return err
}

func (r *rBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	panic("implement me")
}

func (r *rBroker) String() string {
	return "rocketmq"
}
