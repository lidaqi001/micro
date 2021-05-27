package config

const (

	// 每秒钟QPS
	QPS = 100

	// 服务-唱
	SERVICE_SING = "sing"

	// 服务-说
	SERVICE_SPEAK = "speak"

	// 服务-听
	SERVICE_LISTEN = "listen"

	// 服务-异步消息服务-事件
	SERVICE_ASYNC_EVENT = "async_event"

	// 服务-异步消息服务-订阅
	SERVICE_ASYNC_SUBSCRIBER = "async_subscriber"

	// Jaeger 服务端Ip地址
	TRACE_ADDR = "192.168.1.146:6831"

	// 服务注册中心地址
	REGISTER_ADDR = "etcd1:2379"
)
