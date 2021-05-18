package config

// trace端口
//const TRACE_PORT = "6831"

// 每秒钟QPS
//const QPS = 10
//
//// 服务-唱
//const SERVICE_SING = "sing"
//
//// 服务-说
//const SERVICE_SPEAK = "speak"
//
//// 服务-听
//const SERVICE_LISTEN = "listen"
//
//const DEFAULT_TRACE_IP = "192.168.1.146"

const (

	// trace端口
	TRACE_PORT = "6831"
	// 每秒钟QPS
	QPS = 10

	// 服务-唱
	SERVICE_SING = "sing"

	// 服务-说
	SERVICE_SPEAK = "speak"

	// 服务-听
	SERVICE_LISTEN = "listen"

	// Jaeger 服务端Ip地址
	JAEGER_TRACE_IP = "192.168.1.146"
)
