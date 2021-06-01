package rocketmq

import (
	"context"
	"errors"
	"fmt"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/codec/json"
	gerr "github.com/gogap/errors"
	"strings"
	"time"
)

type rBroker struct {
	addrs     []string
	opts      broker.Options
	bopts     *brokerOptions
	client    mq_http_sdk.MQClient
	connected bool
}

type publication struct {
	c       mq_http_sdk.MQConsumer
	topic   string
	message *broker.Message
	err     error
	handles []string
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.message
}

func (p *publication) Ack() error {
	// NextConsumeTime前若不确认消息消费成功，则消息会重复消费
	// 消息句柄有时间戳，同一条消息每次消费拿到的都不一样
	ackerr := p.c.AckMessage(p.handles)
	if ackerr == nil {
		fmt.Printf("Ack ---->\n\t%s\n", p.handles)
	} else {
		// 某些消息的句柄可能超时了会导致确认不成功
		fmt.Println(ackerr)
		for _, errAckItem := range ackerr.(gerr.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem) {
			fmt.Printf("\tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
				errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
	return nil
}

func (p *publication) Error() error {
	return p.err
}

type subscriber struct {
	topic  string
	handle broker.Handler
	opts   broker.SubscribeOptions
}

func (s subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s subscriber) Topic() string {
	return s.topic
}

func (s subscriber) Unsubscribe() error {
	return errors.New("not support unsubscribe~")
}

func init() {
	cmd.DefaultBrokers["rabbitmq"] = NewBroker
}

func NewBroker(opts ...broker.Option) broker.Broker {
	bopts := &brokerOptions{
		instanceId: DefaultKeyValue,
		accessKey:  DefaultKeyValue,
		secretKey:  DefaultKeyValue,
	}
	options := broker.Options{
		// default to json codec
		Codec:   json.Marshaler{},
		Context: context.WithValue(context.Background(), optionsKey, bopts),
	}
	for _, o := range opts {
		o(&options)
	}
	return &rBroker{
		opts:  options,
		bopts: bopts,
	}
}

func (r *rBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

func (r *rBroker) Connect() error {

	switch {
	case len(r.bopts.groupId) == 0:
		return err("groupId")
	case len(r.bopts.endpoint) == 0:
		return err("endpoint")
	case len(r.bopts.accessKey) == 0:
		return err("accessKey")
	case len(r.bopts.secretKey) == 0:
		return err("secretKey")
	case len(r.bopts.instanceId) == 0:
		return err("instanceId")
	}

	client := mq_http_sdk.NewAliyunMQClient(
		r.bopts.endpoint,
		r.bopts.accessKey,
		r.bopts.secretKey,
		"")
	if client != nil {
		r.client = client
		r.connected = true
		return nil
	}

	return errors.New("client is nil!")
}

func (r *rBroker) Disconnect() error {
	r.client = nil
	r.connected = false
	return nil
}

func (r *rBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {

	body, _ := r.opts.Codec.Marshal(m)
	mqProducer := r.client.GetProducer(r.bopts.instanceId, topic)

	var msg mq_http_sdk.PublishMessageRequest
	msg = mq_http_sdk.PublishMessageRequest{
		MessageTag:  "go-micro",          // 消息标签
		MessageBody: string(body),        // 消息内容
		Properties:  map[string]string{}, // 消息属性
	}

	ret, err := mqProducer.PublishMessage(msg)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Publish ---->\n\tMessageId:%s, BodyMD5:%s, \n", ret.MessageId, ret.MessageBodyMD5)
	}
	return nil
}

func (r *rBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber,
	error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
	}
	for _, o := range opts {
		o(&options)
	}

	s := &subscriber{
		topic:  topic,
		handle: handler,
		opts:   options,
	}
	go r.recv(topic, handler, s)
	return s, nil
}

func (r *rBroker) recv(topic string, handler broker.Handler, s *subscriber) {

	consume := r.client.GetConsumer(r.bopts.instanceId, topic, r.bopts.groupId, "")

	for {
		endChan := make(chan int)
		respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
		errChan := make(chan error)
		go func() {
			select {
			case resp := <-respChan:
				{
					// 处理业务逻辑

					var handles []string
					var body string

					fmt.Printf("Consume %d messages---->\n", len(resp.Messages))

					for _, v := range resp.Messages {
						body = v.MessageBody
						handles = append(handles, v.ReceiptHandle)

						var rst *broker.Message
						_ = r.opts.Codec.Unmarshal([]byte(body), &rst)

						p := &publication{
							c:       consume,
							topic:   topic,
							message: rst,
							handles: handles,
						}

						p.err = handler(p)

						if p.err != nil {
							fmt.Printf("p.err:%v", p.err)
							break
						}
						if s.opts.AutoAck {
							_ = p.Ack()
						}

					}
					endChan <- 1
				}
			case err := <-errChan:
				{
					// 没有消息
					if strings.Contains(err.(gerr.ErrCode).Error(), "MessageNotExist") {
						fmt.Println("\nNo new message, continue!")
					} else {
						fmt.Println(err)
						time.Sleep(time.Duration(3) * time.Second)
					}
					endChan <- 1
				}
			case <-time.After(35 * time.Second):
				{
					fmt.Println("Timeout of consumer message ??")
					endChan <- 1
				}
			}
		}()

		// 长轮询消费消息
		// 长轮询表示如果topic没有消息则请求会在服务端挂住3s，3s内如果有消息可以消费则立即返回
		consume.ConsumeMessage(respChan, errChan,
			10, // 一次最多消费3条(最多可设置为16条)
			10, // 长轮询时间3秒（最多可设置为30秒）
		)
		<-endChan
	}
}

func (r *rBroker) String() string {
	return "rocketmq"
}

func (r *rBroker) Options() broker.Options {
	return r.opts
}

func (r *rBroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return "https://ons.console.aliyun.com"
}

func (r *rBroker) Boptions() *brokerOptions {
	return r.bopts
}

func err(param string) error {
	m := fmt.Sprintf("rocketmq：%s param is null!", param)
	return errors.New(m)
}
