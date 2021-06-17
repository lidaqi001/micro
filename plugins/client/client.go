package client

import (
	"context"
	_hystrix "github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/logger"
	"github.com/lidaqi001/micro/plugins/wrapper/breaker/hystrix"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"github.com/opentracing/opentracing-go"
	"io"
)

type client struct {
	opts Options
}

func Create(opts ...Option) (interface{}, error) {

	options := Options{
		Name:     "",
		Init:     nil,
		CallFunc: nil,
		Ctx:      context.Background(),

		Context: context.Background(),
	}

	c := &client{opts: options}

	if err := c.Init(opts...); err != nil {
		return nil, err
	}

	return c.run()
}

func (c *client) Init(opts ...Option) error {

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

	// 设置hystrix默认超时时间（单位：ms）
	_hystrix.DefaultTimeout = 2000
	// hystrix 配置自定义服务（重试、降级、熔断）
	if c.opts.Hystrix != nil {
		hystrix.Configure(c.opts.Hystrix)
	}

	switch {

	case helper.Empty(c.opts.Name):
		return err(NAME_IS_NULL)

	case c.opts.CallFunc == nil:
		return err(CALL_FUNC_IS_NULL)

	}

	return nil
}

func (c *client) run() (interface{}, error) {

	var (
		ctx  = c.opts.Ctx
		name = c.opts.Name

		err   error
		_io   io.Closer
		trace opentracing.Tracer
	)

	// 设置trace client
	trace, _io, err = jaeger.NewTracer(name)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer _io.Close()

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
		micro.WrapClient(traceplugin.NewClientWrapper(trace)),
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
