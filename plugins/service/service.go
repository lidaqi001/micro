package service

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/juju/ratelimit"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/examples/config"
	"sxx-go-micro/plugins/wrapper/service/trace"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

func Create(serviceName string, registerService func(service micro.Service)) {

	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(serviceName, config.TRACE_PORT, "")
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
		// 将服务注册到consul
		micro.Registry(etcd.NewRegistry()),
		// 基于ratelimit 限流
		micro.WrapHandler(
			ratelimiter.NewHandlerWrapper(
				ratelimit.NewBucketWithRate(float64(config.QPS), int64(config.QPS)),
				false),
		),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 链路追踪中间件
		micro.WrapHandler(trace.SpanWrapper),
	)

	// 初始化，会解析命令行参数
	service.Init()

	// 注册处理器，调用服务接口处理请求
	registerService(service)

	// 启动服务
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
