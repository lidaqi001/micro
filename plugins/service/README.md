# service

- 该方法是全局创建服务端的方法

- 实现了
    - 注册发现
    - 链路追踪
    - 服务治理
        - 熔断
        - 重试
        - 降级
        - 限流

- 使用grpc作为传输协议

- 异步消息使用框架默认的 **http broker** 组件

> 但有其致命问题，存储消息的 inbox 参数，最大存储消息只为**64**，多出来的消息不处理
>
> 故实际上不能作为生产环境使用
```go
······

// HTTP Broker is a point to point async broker
type httpBroker struct {

	······

	// offline message inbox
	mtx   sync.RWMutex
	inbox map[string][][]byte
}

······

func (h *httpBroker) saveMessage(topic string, msg []byte) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	// get messages
	c := h.inbox[topic]

	// save message
	c = append(c, msg)

	// max length 64
	if len(c) > 64 {
		c = c[:64]
	}

	// save inbox
	h.inbox[topic] = c
}
```