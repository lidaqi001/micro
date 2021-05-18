package client

import (
	"context"
	"errors"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"log"
	"reflect"
	"sxx-go-micro/common/helper"
	"sxx-go-micro/examples/config"
	hystrix "sxx-go-micro/plugins/wrapper/breaker/hystrix"
	logWrap "sxx-go-micro/plugins/wrapper/client/log"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

// Create params struct
type Params struct {
	ClientName     string
	HystrixService []string
	CallUserFunc   func(micro.Service, context.Context, interface{}) (interface{}, error)
	Ctx            context.Context
	Sp             opentracing.Span
	Input          interface{}
}

// 应用自定义hystrix服务治理的服务列表
var DefaultHystrixService = []string{
	config.SERVICE_SING + ".DemoService.SayHello",
	config.SERVICE_SPEAK + ".DemoService.SayHello",
	config.SERVICE_LISTEN + ".DemoService.SayHello",
}

func Create(params Params) (interface{}, error) {

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
		micro.Registry(etcd.NewRegistry()),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
		// wrap the client
		micro.WrapClient(logWrap.LogWrap),
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

type ClientOp struct {
	Ctx   *gin.Context
	Name  string
	Param Params
}

func GetClient(op ClientOp) (micro.Service, context.Context, error) {

	c := op.Ctx

	var ok bool
	var service interface{}

	if service, ok = c.Get(op.Name); !ok {
		Create(op.Param)
		if service, ok = c.Get("gin"); !ok {
			return nil, nil, errors.New("client create failed!")
		}
		log.Printf("%v", reflect.TypeOf(service))
	}

	cc, _ := c.Get("gin_ctx")
	return service.(micro.Service), cc.(context.Context), nil
}

type cli struct {
	Ctx *gin.Context
}

type Cli interface {
	Create(Params)
	GetClient(ClientOp) (micro.Service, context.Context, error)
}
