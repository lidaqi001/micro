# 弃用

Broker驱动 =》 阿里云【消息队列 for Apache RocketMQ】


### 问题：

- 多subscribe订阅同一TOPIC时，tag过滤无效（原因不明，sdk中代码也基本无注释）
    
- 使用指定broker驱动后，micro.WrapSubscriber() 函数不生效
（后期在深入源码看一下）


- ！！！http协议目前不支持广播消费模式！！！
    - [广播模式/集群模式-介绍](http://www.baidu.com)

这是阿里云客服的回复：

```text
消息队列 RocketMQ http sdk 如何实现广播消费模式
您好，因为产品限制， http协议目前不支持广播消费模式；

TCP协议可以支持广播消费模式详情请参考：Java订阅消息、.NET订阅消息、C/C++：订阅消息。

你也可以考虑多个group去消费同一个topic的消息，模拟广播消费模式。
```
