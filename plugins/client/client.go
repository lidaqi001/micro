package client

import (
	"context"
	"errors"
	hystrix2 "github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/logger"
	"github.com/lidaqi001/micro/plugins/wrapper/breaker/hystrix"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
)

type client struct {
	opts Options
}

func Create(opts ...Option) (interface{}, error) {
	options := Options{
		Name:     "",
		Init:     nil,
		CallFunc: nil,
		Context:  context.Background(),
	}
	c := &client{opts: options}
	return c.Init(opts...)
}

func (c client) Init(opts ...Option) (interface{}, error) {

	for _, o := range opts {
		o(&c.opts)
	}

	if val, ok := c.opts.Context.Value(nameKey{}).(string); ok {
		c.opts.Name = val
	}
	if val, ok := c.opts.Context.Value(hystrixKey{}).([]string); ok {
		c.opts.Hystrix = val
	}
	if val, ok := c.opts.Context.Value(ctxKey{}).(context.Context); ok {
		c.opts.Ctx = val
	}
	if val, ok := c.opts.Context.Value(inputKey{}).(interface{}); ok {
		c.opts.Input = val
	}
	if val, ok := c.opts.Context.Value(initKey{}).([]micro.Option); ok {
		c.opts.Init = val
	}
	if val, ok := c.opts.Context.Value(callFuncKey{}).(func(CallFuncParams) (interface{}, error)); ok {
		c.opts.CallFunc = val
	}

	switch {

	case helper.Empty(c.opts.Name):
		err := errors.New(NAME_IS_NULL)
		logger.Error(err)
		return nil, err

	case c.opts.CallFunc == nil:
		err := errors.New(CALL_FUNC_IS_NULL)
		logger.Error(err)
		return nil, err
	}

	return c.run()
}

func (c client) run() (interface{}, error) {

	var (
		ctx  = c.opts.Ctx
		name = c.opts.Name
	)

	// 设置hystrix默认超时时间（单位：ms）
	hystrix2.DefaultTimeout = 2000

	if c.opts.Hystrix != nil {
		// hystrix 配置自定义服务（重试、降级、熔断）
		hystrix.Configure(c.opts.Hystrix)
	}

	//if ctx == nil || ctx == context.Background() {
	//	// 当 ctx==nil || ctx==context.Background() 时
	//	// 初始化上下文和span
	//	_, ctx = jaeger.GetTraceClientCtxAndSpan()
	//}

	// 设置trace server地址
	t, io, err := jaeger.NewTracer(name)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer io.Close()

	// 创建一个新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		// 客户端名称
		micro.Name(name),
		// 服务发现
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(helper.GetRegistryAddress()),
		)),
		// 使用 hystrix 实现服务治理
		micro.WrapClient(hystrix.NewClientWrapper()),
		// 链路追踪客户端
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
		// 自定义客户端中间件
		//micro.WrapClient(log.LogWrap),
	)

	// 初始化
	service.Init(c.opts.Init...)

	// 执行客户端闭包，调用相应服务
	return c.opts.CallFunc(CallFuncParams{
		Ctx:     ctx,
		Service: service,
		Input:   c.opts.Input,
	})

}
