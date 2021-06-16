# asyncRocketmq（不可用）

- 该方法是在 ***plugins/service/service.go*** 的基础上

- 使用了自己封装的 **rocketmq（阿里云）** 驱动 ***plugins/broker/rocketmq***

- 使异步消息在更稳定的消息队列中运行

### 问题：
    
```text
使用指定broker驱动后，micro.WrapSubscriber() 函数不生效
（后期在深入源码看一下）
```