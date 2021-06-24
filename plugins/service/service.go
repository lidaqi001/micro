package service

import (
	"context"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/juju/ratelimit"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/broker/rabbitmq"
	"github.com/lidaqi001/micro/plugins/logger"
	"github.com/lidaqi001/micro/plugins/wrapper/service/trace"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"github.com/opentracing/opentracing-go"
	"time"
)

type service struct {
	opts Options
}

var (
	serverInit = []server.Option{}
)

func Create(opts ...Option) error {
	options := Options{
		Name:      "",
		Advertise: "",
		Init:      nil,
		CallFunc:  nil,
		Context:   context.Background(),
		Rabbitmq:  false,
	}

	s := &service{opts: options}

	if err := s.Init(opts...); err != nil {
		return err
	}

	return s.run()
}

func (s *service) Init(opts ...Option) error {

	for _, o := range opts {
		o(&s.opts)
	}

	if val, ok := s.opts.Context.Value(advertiseKey{}).(string); ok {
		s.opts.Advertise = val
	}
	if val, ok := s.opts.Context.Value(serviceNameKey{}).(string); ok {
		s.opts.Name = val
	}
	if val, ok := s.opts.Context.Value(initKey{}).([]micro.Option); ok {
		s.opts.Init = val
	}
	if val, ok := s.opts.Context.Value(callFuncKey{}).(func(micro.Service)); ok {
		s.opts.CallFunc = val
	}
	if val, ok := s.opts.Context.Value(Rabbitmq{}).(bool); ok && val {
		// 设置rabbitmq地址
		rabbitmq.DefaultRabbitURL = helper.GetConfig("RABBITMQ_ADDR", config.RABBITMQ_ADDR)
		s.opts.Init = append(s.opts.Init, micro.Broker(
			// 设置 rabbitmq 为 broker 驱动
			rabbitmq.NewBroker(
				// 设置：Exchange 为持久化
				// If this option is not set, the exchange will be deleted when rabbitmq restarts
				rabbitmq.DurableExchange(),
				// 设置：订阅时创建持久化队列
				rabbitmq.PrefetchGlobal(),
			),
		),
		)
	}

	switch {

	case helper.Empty(s.opts.Name):
		return err(SERVICE_NAME_IS_NULL)

	case s.opts.CallFunc == nil:
		return err(CALL_FUNC_IS_NULL)

	case !helper.Empty(s.opts.Advertise):
		serverInit = append(serverInit, server.Advertise(s.opts.Advertise))
		// fallthrough 关键字，继续向下匹配case
		//fallthrough

	}

	return nil
}

func (s *service) run() error {
	name := s.opts.Name

	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(name)
	if err != nil {
		logger.Error(err)
		return err
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
		micro.Name(name),
		// 服务注册
		micro.RegisterTTL(time.Second*30),      // 注册存活时间30s
		micro.RegisterInterval(time.Second*20), // 刷新服务存活时间的间隔时间（保持注册存活时间）
		micro.Registry(etcd.NewRegistry(registry.Addrs(helper.GetRegistryAddress()))),
		// wrap handler
		micro.WrapHandler(
			// 基于ratelimit 限流
			ratelimiter.NewHandlerWrapper(ratelimit.NewBucketWithRate(helper.GetQPS()), false),
		),
		micro.WrapHandler(
			// 基于 jaeger 采集追踪数据
			// handler 调用服务-链路追踪
			traceplugin.NewHandlerWrapper(t),
			trace.SpanWrapper,
		),
		// wrap subscriber
		// subscriber 消息服务（异步事件/订阅）-链路追踪
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(t),
			trace.SubWrapper,
		),
	)

	// 初始化，会解析命令行参数
	service.Init(s.opts.Init...)

	// 服务初始化
	service.Server().Init(serverInit...)

	// 注册处理器，调用服务接口处理请求
	s.opts.CallFunc(service)

	// 启动服务
	if err := service.Run(); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
