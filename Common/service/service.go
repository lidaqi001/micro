package service

import (
	"context"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/juju/ratelimit"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/Common/config"
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

	// 限流
	bucket := ratelimit.NewBucketWithRate(float64(config.QPS), int64(config.QPS))

	// 创建新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Server(grpc.NewServer()),
		// 服务名称
		micro.Name(serviceName),
		// 将服务注册到consul
		micro.Registry(consul.NewRegistry()),
		// 基于ratelimit 限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(bucket, false)),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
		// wrap the handler
		micro.WrapHandler(spanWrapper),
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

// spanWrapper is a handler wrapper
func spanWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		//log.Printf("[wrapper] server request: %v", req.Endpoint())
		//log.Printf("[wrapper] server request params: %v", req.Body())

		sp := jaeger.NewSpan(ctx)
		// Trace：请求service前打印请求tag
		sp.SetReq(req.Body())

		// 执行 service 注册函数
		err := fn(ctx, req, rsp)

		// Trace：执行函数后打印 返回值/错误 tag
		sp.SetRes(rsp, err)

		//log.Printf("[wrapper] server rsp: %v", rsp)
		return err
	}
}
