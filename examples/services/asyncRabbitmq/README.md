# Broker Driver：rabbitmq

- 可用

- 问题
    - rabbitmq消息持久化option
        - rabbitmq.DeliveryMode(amqp.Persistent))
        - 两种方案如何选择？
        - 选第一种：不能设置rabbitmq消息持久化
        - 选第二种，不能使用**micro.WrapSubscriber()** 中间件
            - 第一种：
                - 发布：micro.NewEvent().Publish()
                - 订阅：micro.RegisterSubscriber()
                    - 这一套组合不能设置rabbitmq的消息持久化
                    - 因为返回类型是 **broker.PublishOption**
            - 第二种：
                - 发布：**s.Service.Options().Broker.Publish()** 
                - 订阅：**s.Service.Options().Broker.Subscribe()** 
                    - 这一套组合可以设置消息持久化option
                    - 但是不能使用 **micro.WrapSubscriber()** 中间件