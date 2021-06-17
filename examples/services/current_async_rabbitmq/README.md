# Broker Driver：rabbitmq

- 可用

- 问题
    - rabbitmq消息持久化option（rabbitmq.DeliveryMode(amqp.Persistent))）
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
        
        - 最终解决方案：
            - 重写了rabbitmq驱动（plugins/broker/rabbitmq）
            - 仍然使用 micro 发布/订阅，封装组件（plugins/rabbitmqPack）
                - 发布：micro.NewEvent().Publish()
                - 订阅：micro.RegisterSubscriber()