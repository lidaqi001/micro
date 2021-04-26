package client

import (
	"context"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	log2 "github.com/asim/go-micro/v3/util/log"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/Common/wrapper/breaker/hystrix"
	"sxx-go-micro/trace"
)

func CreateClient(
	clientName string,
	callUserFunc func(service micro.Service, ctx context.Context) (interface{}, interface{}, error),
	ctx context.Context,
	sp opentracing.Span) (interface{}, error) {

	// 这一行封装出去有问题
	//t, sp, parentContext := traceGetConfig(clientName, ctx, sp)

	// 设置trace server地址
	/*
		1/ 从外部设置env		PS：traceIp := os.Getenv("MICRO_TRACE_SERVER")
		2/ 代码中设置env		PS：traceIp := os.Setenv("MICRO_TRACE_IP", "192.168.1.146")

		参数传递：
			t, io, err := trace.NewTracer("service.trace", traceServer, traceIp)
	*/
	t, io, err := trace.NewTracer(clientName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// 获取父级上下文
	// 如果没有传入上下文和span，则创建
	var parentContext context.Context
	if ctx == nil || sp == nil {
		// 创建空的上下文, 生成追踪 span
		// 注入 opentracing textmap 到空的上下文用于追踪
		sp, parentContext = trace.GetTraceClientCtx()
	} else {
		parentContext = opentracing.ContextWithSpan(ctx, sp)
	}





	// hystrix 配置
	hystrix.Configure([]string{"sing.DemoService.SayHello"})

	// 创建一个新的服务
	service := micro.NewService(
		micro.Name(clientName),
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		// 客户端从consul中发现服务
		micro.Registry(consul.NewRegistry()),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
	)
	// 初始化
	service.Init()

	// 执行客户端调用
	req, resp, err := callUserFunc(service, parentContext)

	// 设置 trace tag
	//traceSetTag(sp, req, resp, err)
	sp.SetTag("req", req)

	if err != nil {
		// 记录错误
		sp.SetTag("err", err)
		//log.Println("服务调用失败：%v", err)
		//log.Fatalf("服务调用失败：%v", err)
		log2.Warnf("服务调用失败：%v", err)
		return resp, err
	}

	// 记录响应
	sp.SetTag("resp", resp)

	return resp, err
}

func traceSetTag(sp opentracing.Span, args ...interface{}) {

	// 记录请求
	sp.SetTag("req", args[0])

	err := args[2].(error)
	if err != nil {
		// 记录错误
		sp.SetTag("err", err)
		log.Println("服务调用失败：%v", err)
		//log.Fatalf("服务调用失败：%v", err)
		return
	}

	// 记录响应
	sp.SetTag("resp", args[1])
}

func traceGetConfig(
	clientName string,
	ctx context.Context,
	sp opentracing.Span) (opentracing.Tracer, opentracing.Span, context.Context) {

	// 设置trace server地址
	/*
		1/ 从外部设置env		PS：traceIp := os.Getenv("MICRO_TRACE_SERVER")
		2/ 代码中设置env		PS：traceIp := os.Setenv("MICRO_TRACE_IP", "192.168.1.146")

		参数传递：
			t, io, err := trace.NewTracer("service.trace", traceServer, traceIp)
	*/
	t, io, err := trace.NewTracer(clientName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// 获取父级上下文
	// 如果没有传入上下文和span，则创建
	var parentContext context.Context
	if ctx == nil || sp == nil {
		// 创建空的上下文, 生成追踪 span
		// 注入 opentracing textmap 到空的上下文用于追踪
		sp, parentContext = trace.GetTraceClientCtx()
	} else {
		parentContext = opentracing.ContextWithSpan(ctx, sp)
	}

	return t, sp, parentContext
}
