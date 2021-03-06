package service

import (
	"context"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
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
		CallFunc:   func(m micro.Service) {},
		Context:    context.Background(),
		RpcAddr:    DEFAULT_RPC_ADDR,
		HttpAddr:   DEFAULT_HTTP_ADDR,
		Rabbitmq:   true,
		Server:     grpc.NewServer(),
		ServerType: RPC,
		BindRoute:  DEFAULT_BIND_ROUTE,
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

	// set rabbitmq for broker driver
	if val, ok := s.opts.Context.Value(rabbitmqKey{}).(bool); ok {
		s.opts.Rabbitmq = val
	}
	if val, ok := s.opts.Context.Value(addressKey{}).(string); ok {
		s.opts.RpcAddr = val
		s.opts.HttpAddr = val
	}
	if val, ok := s.opts.Context.Value(serverTypeKey{}).(Mode); ok {
		s.opts.ServerType = val
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
	if val, ok := s.opts.Context.Value(bindRouteKey{}).(func(*gin.Engine)); ok {
		s.opts.BindRoute = val
	}
	if val, ok := s.opts.Context.Value(callFuncKey{}).(func(micro.Service)); ok {
		s.opts.CallFunc = val
	}

	switch {

	case helper.Empty(s.opts.Name):
		return sErr(SERVICE_NAME_IS_NULL)

	//case s.opts.CallFunc == nil:
	//	return sErr(CALL_FUNC_IS_NULL)

	case !helper.Empty(s.opts.Advertise):
		serverInit = append(serverInit, server.Advertise(s.opts.Advertise))
		// fallthrough ??????????????????????????????case
		//fallthrough

	}

	// set rabbitmq
	if s.opts.Rabbitmq == true {
		rabbitmq.DefaultRabbitURL = helper.GetRabbitmqAddress()
		s.opts.Init = append(
			s.opts.Init, micro.Broker(
				// ?????? rabbitmq ??? broker ??????
				rabbitmq.NewBroker(
					// ?????????Exchange ????????????
					// If this option is not set, the exchange will be deleted when rabbitmq restarts
					rabbitmq.DurableExchange(),
					// ???????????????????????????????????????
					rabbitmq.PrefetchGlobal(),
				),
			),
		)
	}

	switch s.opts.ServerType {
	case HTTP:
		err := s.lnitializeHttp()
		if err != nil {
			return sErr(err)
		}
	case RPC:
		s.opts.Init = append(s.opts.Init, micro.Address(s.opts.RpcAddr))
	}

	return nil
}

func (s *service) lnitializeHttp() error {

	srv := httpServer.NewServer(
		server.Name(s.opts.Name+"_http"),
		server.Address(s.opts.HttpAddr),
	)

	// set gin mode
	v, err := helper.IsOpenDebug()
	if err == nil {
		if v {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		return err
	}

	router := gin.Default()

	// bind route for http server
	s.opts.BindRoute(router)

	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		logger.Error(err)
		return err
	}

	s.opts.Server = srv

	return nil
}

func (s *service) run() error {
	name := s.opts.Name

	// ???????????????????????????
	t, io, err := jaeger.NewTracer(name)
	if err != nil {
		return sErr(err)
	}
	defer io.Close()
	// ??????????????? tracer
	// PS????????? service ??????
	opentracing.SetGlobalTracer(t)

	// ??????????????????
	service := micro.NewService(
		// ??????grpc??????
		micro.Server(s.opts.Server),
		// ????????????
		micro.Name(name),
		// ????????????
		micro.RegisterTTL(time.Second*30),      // ??????????????????30s
		micro.RegisterInterval(time.Second*20), // ?????????????????????????????????????????????????????????????????????
		micro.Registry(helper.GetRegistry()),
		// wrap handler
		micro.WrapHandler(
			// ??????ratelimit ??????
			ratelimiter.NewHandlerWrapper(ratelimit.NewBucketWithRate(helper.GetQPS()), false),

			// ?????? jaeger ??????????????????
			// handler ????????????-????????????
			traceplugin.NewHandlerWrapper(t),
			trace.SpanWrapper,
		),
		// wrap subscriber
		// subscriber ???????????????????????????/?????????-????????????
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(t),
			trace.SubWrapper,
		),
	)

	// ????????????????????????????????????
	service.Init(s.opts.Init...)

	// ???????????????
	err = service.Server().Init(serverInit...)
	if err != nil {
		return sErr(err)
	}

	// ????????????????????????????????????????????????
	s.opts.CallFunc(service)

	// ????????????
	if err := service.Run(); err != nil {
		return sErr(err)
	}

	return nil
}
