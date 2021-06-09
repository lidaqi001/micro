package service

import (
	"context"
	"errors"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/juju/ratelimit"
	"github.com/lidaqi001/micro/common"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	log "github.com/lidaqi001/micro/plugins/logger"
	"github.com/lidaqi001/micro/plugins/wrapper/service/trace"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"github.com/opentracing/opentracing-go"
	"os"
)

type service struct {
	opts Options
}

func Create(opts ...Option) {
	options := Options{
		Context:     context.Background(),
		ServiceName: "",
		Init:        nil,
		CallFunc:    nil,
	}
	service := &service{opts: options}
	service.Init(opts...)
}

func (s *service) Init(opts ...Option) {

	logger.DefaultLogger = log.NewLogger(
		// 日志目录
		log.OutputFilePath(config.LOG_DEFAULT_SERVICE),
		// 日志根目录
		log.OutputRootPath(config.LOG_ROOT),
	)
	//common.SetDefaultLoggerForZerolog(config.LOG_DEFAULT_SERVICE)
	for _, o := range opts {
		o(&s.opts)
	}
	if name, ok := s.opts.Context.Value(serviceNameKey{}).(string); ok {
		s.opts.ServiceName = name
	}
	if init, ok := s.opts.Context.Value(initKey{}).([]micro.Option); ok {
		s.opts.Init = init
	}
	if fn, ok := s.opts.Context.Value(callFuncKey{}).(func(micro.Service)); ok {
		s.opts.CallFunc = fn
	}

	switch {
	case helper.Empty(s.opts.ServiceName):
		//logger.Fatal("The serviceName cannot be empty~")
		logger.Error("The serviceName cannot be empty~")
		logger.Debug("The serviceName cannot be empty~")
		os.Exit(1)
	case s.opts.CallFunc == nil:
		logger.Fatal(errors.New("The CallFunc cannot be empty~"))
	}
}

func (s *service) run() {

	common.SetDefaultLoggerForZerolog(config.LOG_DEFAULT_SERVICE)

	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(s.opts.ServiceName)
	if err != nil {
		logger.Fatal(err)
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
		micro.Name(s.opts.ServiceName),
		// 服务注册
		micro.Registry(etcd.NewRegistry(registry.Addrs(helper.GetRegistryAddress()))),
		// wrap handler
		micro.WrapHandler(
			// 基于ratelimit 限流
			ratelimiter.NewHandlerWrapper(
				ratelimit.NewBucketWithRate(helper.GetQPS()), false,
			),
		),
		micro.WrapHandler(
			// 基于 jaeger 采集追踪数据
			// handler 调用服务-链路追踪
			traceplugin.NewHandlerWrapper(t),
			trace.SpanWrapper,
		),
		// wrap subscriber
		// subscriber 消息服务（异步事件/订阅）-链路追踪
		// ！！！目前对于自定义的驱动不生效！！！
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(t),
			trace.SubWrapper,
		),
	)

	// 初始化，会解析命令行参数
	service.Init(s.opts.Init...)

	// 将broker设置为不限制ip，默认为127.0.0.1
	//_ = service.Options().Broker.Init(
	//	broker.Addrs("0.0.0.0:"),
	//)

	// 注册处理器，调用服务接口处理请求
	s.opts.CallFunc(service)

	// 启动服务
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
