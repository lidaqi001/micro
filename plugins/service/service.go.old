package service

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/juju/ratelimit"
	"github.com/lidaqi001/micro/common"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/wrapper/service/trace"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"github.com/opentracing/opentracing-go"
	"log"
)

func Create(serviceName string, registerService func(service micro.Service), opts ...micro.Option) {
	common.SetDefaultLoggerForZerolog(config.LOG_DEFAULT_SERVICE)

	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	// 设置为全局 tracer
	// PS：只在 service 声明
	opentracing.SetGlobalTracer(t)

	// 创建新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Server(grpc.NewServer()),
		// 服务名称
		micro.Name(serviceName),
		// 服务注册
		micro.Registry(etcd.NewRegistry(registry.Addrs(helper.GetRegistryAddress()))),
		// wrap handler
		micro.WrapHandler(
			// 基于ratelimit 限流
			ratelimiter.NewHandlerWrapper(
				ratelimit.NewBucketWithRate(helper.GetQPS()), false,
			),
		),
		micro.WrapHandler(
			// 基于 jaeger 采集追踪数据
			// handler 调用服务-链路追踪
			traceplugin.NewHandlerWrapper(t),
			trace.SpanWrapper,
		),
		// wrap subscriber
		// subscriber 消息服务（异步事件/订阅）-链路追踪
		// ！！！目前对于自定义的驱动不生效！！！
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(t),
			trace.SubWrapper,
		),
	)

	// 初始化，会解析命令行参数
	service.Init(opts...)

	// 将broker设置为不限制ip，默认为127.0.0.1
	//_ = service.Options().Broker.Init(
	//	broker.Addrs("0.0.0.0:"),
	//)

	// 注册处理器，调用服务接口处理请求
	registerService(service)

	// 启动服务
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
