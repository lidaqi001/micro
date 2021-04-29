package client

import (
	"context"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/plugins/wrapper/breaker/hystrix"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

// 此结构体用以改造Create函数，非必须参数改用Params结构体存储
// 未完成
type Params struct {
	Ctx            context.Context
	Sp             opentracing.Span
	HystrixService string
}

func Create(
	clientName string,
	callUserFunc func(service micro.Service, ctx context.Context) (interface{}, interface{}, error),
	ctx context.Context,
	sp opentracing.Span,
	hystrixService []string) (interface{}, error) {

	// 设置trace server地址
	t, io, err := jaeger.NewTracer(clientName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// hystrix 配置（重试、降级、熔断）
	hystrix.Configure(hystrixService)

	// 创建一个新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		// 客户端名称
		micro.Name(clientName),
		// 客户端从consul中发现服务
		micro.Registry(consul.NewRegistry()),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
	)

	// 初始化
	service.Init()

	// 当ctx || sp 为空时
	if ctx == nil || sp == nil {
		sp, ctx = jaeger.GetTraceClientCtxAndSpan()
	}

	// 执行客户端闭包，调用相应服务
	req, resp, err := callUserFunc(service, ctx)

	// 记录请求
	jaeger.SpanSetRequest(sp, req)
	defer func() {
		// 记录错误
		// 记录响应
		jaeger.SpanSetResClient(sp, resp, err)
		if err != nil {
			log.Printf("服务调用失败：%v", err)
		}
	}()

	return resp, err
}
