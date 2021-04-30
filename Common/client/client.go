package client

import (
	"context"
	"errors"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/Common/helper"
	"sxx-go-micro/plugins/wrapper/breaker/hystrix"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

// log wrapper logs every time a request is made
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	log.Printf("[wrapper] client request service: %s method: %s\n", req.Service(), req.Endpoint())

	// 请求服务
	err := l.Client.Call(ctx, req, rsp)

	log.Printf("[wrapper] client rsp: %v\n", rsp)
	return err
}

// Implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

// Create params struct
type Params struct {
	ClientName     string
	HystrixService []string
	CallUserFunc   func(service micro.Service, ctx context.Context) (interface{}, error)
	Ctx            context.Context
	Sp             opentracing.Span
}

func Create(params *Params) (interface{}, error) {

	err := verifyParams(*params)
	if err != nil {
		return nil, err
	}

	if params.HystrixService != nil {
		// hystrix 配置（重试、降级、熔断）
		hystrix.Configure(params.HystrixService)
	}

	sp, ctx := params.Sp, params.Ctx
	if sp == nil || ctx == nil {
		// 当ctx || sp 为空时
		// 初始化上下文和span
		sp, ctx = jaeger.GetTraceClientCtxAndSpan()
	}

	// 设置trace server地址
	t, io, err := jaeger.NewTracer(params.ClientName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// 创建一个新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		// 客户端名称
		micro.Name(params.ClientName),
		// 客户端从consul中发现服务
		micro.Registry(consul.NewRegistry()),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
		// wrap the client
		micro.WrapClient(logWrap),
	)

	// 初始化
	service.Init()

	// 执行客户端闭包，调用相应服务
	return params.CallUserFunc(service, ctx)
}

func verifyParams(params Params) error {
	switch {
	case helper.Empty(params.ClientName):
		return errors.New("clientName is empty!")
	case params.CallUserFunc == nil:
		return errors.New("CallUserFunc is nil!")
	}
	return nil
}
