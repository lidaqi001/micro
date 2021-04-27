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
	"sxx-go-micro/plugins/wrapper/breaker/hystrix"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

func CreateClient(
	clientName string,
	callUserFunc func(service micro.Service, ctx context.Context) (interface{}, interface{}, error),
	ctx context.Context,
	sp opentracing.Span,
	hystrixService []string) (interface{}, error) {

	// 当ctx || sp 为空时
	if ctx == nil || sp == nil {
		sp, ctx = jaeger.GetTraceClientCtxAndSpan()
	}

	// 设置trace server地址
	t, io, err := jaeger.NewTracer(clientName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// hystrix 配置
	hystrix.Configure(hystrixService)

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

	// 执行客户端闭包，调用相应服务
	req, resp, err := callUserFunc(service, ctx)

	// 记录请求
	sp.SetTag("req", req)
	defer func() {
		if err != nil {
			// 记录错误
			sp.SetTag("err", err)
			log2.Warnf("服务调用失败：%v", err)
		}
		// 记录响应
		sp.SetTag("resp", resp)
	}()

	return resp, err
}
