# rabbitmq broker driver

> 基于 go-micro 原 rabbitmq 驱动组件修改

> 用以解决 micro.NewEvent 不能设置rabbitmq持久化消息的配置
    
   - rabbitmq.DeliveryMode(amqp.Persistent)
   
> 修改

- options.go
```go

//type deliveryMode struct{}
//type priorityKey struct{}

// 改为公共参数，可在外部设置

type DeliveryMode struct{}
type Priority struct{}
```

- rabbitmq.go
```go

func (r *rbroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
······
	if options.Context != nil {
		if value, ok := options.Context.Value(DeliveryMode{}).(uint8); ok {
			fmt.Println("DeliveryMode:", value)
			m.DeliveryMode = value
		}

		if value, ok := options.Context.Value(Priority{}).(uint8); ok {
			fmt.Println("Priority:", value)
			m.Priority = value
		}
	}
······
}
```

- examples/services/asyncRabbitmq/event/handler/publish.go
```text

func (s *DemoServiceHandler) publishSayHello(req string) error {

	// 设置DeliveryMode
	ctx := context.WithValue(context.Background(), rabbitmq.DeliveryMode{}, amqp.Persistent)
    // 设置Priority
	ctx = context.WithValue(ctx, rabbitmq.Priority{}, 0)

	// 发布消息
	p := micro.NewEvent(config.EVENT_B, s.Service.Client())
	if err := p.Publish(context.TODO(), &event.DemoEvent{
		City:        req,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	}, 
            // 设置 PublishContext
            client.PublishContext(ctx),
	); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}
```