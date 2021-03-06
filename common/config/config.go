package config

const (

	/*******************************************************
					rocketmq（阿里云）配置
	*******************************************************/

	ROCKETMQ_TOPIC_DEFAULT = "go_micro"
	ROCKETMQ_INSTANCE_ID   = ""
	ROCKETMQ_ACCESS_KEY    = ""
	ROCKETMQ_SECRET_KEY    = ""
	ROCKETMQ_ENDPOINT      = ""
	ROCKETMQ_GROUP_ID      = ""

	/*******************************************************
				            日志目录
	*******************************************************/
	LOG_ROOT    = "/microRuntime/"
	LOG_DEFAULT = "log/"

	/*******************************************************
							限流配置
	*******************************************************/

	// 每秒钟QPS
	QPS = "100"

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

	/*******************************************************
						链路追踪配置
	*******************************************************/

	// Jaeger 服务端Ip地址
	TRACE_ADDR = "127.0.0.1:6831"

	/*******************************************************
					服务注册发现中心（etcd）配置
	*******************************************************/

	REGISTRY_ADDR = "127.0.0.1:2379"

	/*******************************************************
						RabbitMq
	*******************************************************/

	RABBITMQ_ADDR = "amqp://guest:guest@127.0.0.1:5672"

	/*******************************************************
							事件
	*******************************************************/
	EVENT_A = "callSing"
	EVENT_B = "singEven2t"

	/*******************************************************
				            队列
	*******************************************************/
	QUEUE_A = "A"
	QUEUE_B = "B"
	QUEUE_C = "C"

	/*******************************************************
				            mysql
	*******************************************************/
	M_USER     = ""
	M_PASSWORD = ""
	M_HOST     = "127.0.0.1:3306"
	M_DBNAME   = "dbname"
)
