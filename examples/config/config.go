package config

const (

	/*******************************************************
					rocketmq（阿里云）配置
	*******************************************************/

	ROCKETMQ_GROUP_ID      = ""
	ROCKETMQ_ENDPOINT      = ""
	ROCKETMQ_ACCESS_KEY    = ""
	ROCKETMQ_SECRET_KEY    = ""
	ROCKETMQ_INSTANCE_ID   = ""
	ROCKETMQ_TOPIC_DEFAULT = "go_micro"

	/*******************************************************
						限流配置
	*******************************************************/

	// 每秒钟QPS
	QPS = 100

	/*******************************************************
						服务配置
	*******************************************************/

	// 唱
	SERVICE_SING = "sing"

	// 说
	SERVICE_SPEAK = "speak"

	// 听
	SERVICE_LISTEN = "listen"

	// 异步消息服务-事件
	SERVICE_ASYNC_EVENT = "async_event"

	// 异步消息服务-订阅
	SERVICE_ASYNC_SUBSCRIBER = "async_subscriber"

	// 异步消息服务-事件-rocketmq
	SERVICE_ASYNC_EVENT_ROCKETMQ = "async_event_rocketmq"

	// 异步消息服务-订阅-rocketmq
	SERVICE_ASYNC_SUBSCRIBER_ROCKETMQ = "async_subscriber_rocketmq"

	/*******************************************************
						链路追踪配置
	*******************************************************/

	// Jaeger 服务端Ip地址
	TRACE_ADDR = "127.0.0.1:6831"

	/*******************************************************
					服务注册中心（etcd）配置
	*******************************************************/

	// 服务注册中心地址
	REGISTRY_ADDR = "127.0.0.1:2379"
)
