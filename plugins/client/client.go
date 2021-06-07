package client

import (
	"context"
	"errors"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/lidaqi001/micro/common"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	hystrix "github.com/lidaqi001/micro/plugins/wrapper/breaker/hystrix"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"log"
)

// Create params struct
type Params struct {
	ClientName     string
	HystrixService []string
	CallUserFunc   func(micro.Service, context.Context, interface{}) (interface{}, error)
	Ctx            context.Context
	Input          interface{}
}

// 应用自定义hystrix服务治理的服务列表
var DefaultHystrixService = []string{
	//config.SERVICE_SING + ".DemoService.SayHello",
	//config.SERVICE_SPEAK + ".DemoService.SayHello",
	//config.SERVICE_LISTEN + ".DemoService.SayHello",
}

func Create(params Params) (interface{}, error) {
	common.SetDefaultLoggerForZerolog(config.LOG_DEFAULT_CLIENT)

	err := verifyParams(params)
	if err != nil {
		return nil, err
	}

	if params.HystrixService != nil {
		// hystrix 配置（重试、降级、熔断）
		hystrix.Configure(params.HystrixService)
	} else {
		hystrix.Configure(DefaultHystrixService)
	}

	ctx := params.Ctx
	if ctx == nil {
		// 当ctx || sp 为空时
		// 初始化上下文和span
		_, ctx = jaeger.GetTraceClientCtxAndSpan()
	}

	// 设置trace server地址
	t, io, err := jaeger.NewTracer(params.ClientName)
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
		// 服务发现
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(helper.GetRegistryAddress()),
		)),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
		// wrap the client
		//micro.WrapClient(logWrap.LogWrap),
	)

	// 初始化
	service.Init()

	// 执行客户端闭包，调用相应服务
	return params.CallUserFunc(service, ctx, params.Input)
}

func verifyParams(params Params) error {
	switch {
	case helper.Empty(params.ClientName):
		return errors.New("clientName can't be empty!")
	case params.CallUserFunc == nil:
		return errors.New("CallUserFunc can't be nil!")
	}
	return nil
}
