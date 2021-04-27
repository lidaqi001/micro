package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"reflect"
	"sxx-go-micro/Common/client"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
	"sxx-go-micro/proto"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *proto.UserRequest, *proto.DemoResponse) error {
	panic("implement me")
}

//func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {
//	// 从微服务上下文中获取追踪信息
//	md, ok := metadata.FromContext(ctx)
//	if !ok {
//		md = make(map[string]string)
//	}
//	var sp opentracing.Span
//	wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
//	// 创建新的 Span 并将其绑定到微服务上下文
//	sp = opentracing.StartSpan("SayHello", opentracing.ChildOf(wireContext))
//	// 记录请求
//	sp.SetTag("req", req)
//	var err error
//	err = nil
//	defer func() {
//		if err!=nil{
//			sp.SetTag("err", err)
//		}
//		// 记录响应
//			sp.SetTag("res", rsp)
//		// 在函数返回 stop span 之前，统计函数执行时间
//		sp.Finish()
//	}()
//
//	clientName := "client.2"
//	// 设置trace server地址
//	t, io, err := trace.NewTracer(clientName, config.TRACE_PORT, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer io.Close()
//
//
//	// hystrix 配置
//	//hystrix.Configure(hystrixService)
//	//hystrixGo.DefaultTimeout=5
//
//	// 创建一个新的服务
//	service := micro.NewService(
//		micro.Name(clientName),
//		// 使用grpc协议
//		micro.Client(grpc.NewClient()),
//		// 客户端从consul中发现服务
//		micro.Registry(consul.NewRegistry()),
//		// 使用 hystrix 实现服务治理
//		//micro.WrapClient(hystrix.NewClientWrapper()),
//		// 链路追踪客户端
//		micro.WrapClient(traceplugin.NewClientWrapper(t)),
//	)
//	// 初始化
//	service.Init()
//
//	// 执行客户端调用
//	cli := proto.NewDemoService(config.SERVICE_SPEAK, service.Client())
//	reqq := &proto.DemoRequest{Name: "学院君"}
//	resp, err := cli.SayHello(ctx, reqq)
//
//	rsp.Text = resp.Text
//
//	return nil
//}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	sp := jaeger.GetTraceServiceSpan(&ctx, req, nil)

	// 记录响应
	// 在函数返回 stop span 之前，统计函数执行时间
	defer jaeger.SpanSetResponse(sp, rsp)

	// 调用service2
	res, _ := client.CreateClient(
		"client.2",
		func(service micro.Service, ctx context.Context) (interface{}, interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_SPEAK, service.Client())
			req := &proto.DemoRequest{Name: "学院君"}
			resp, err := cli.SayHello(ctx, req)
			return req, resp, err
		}, ctx, sp,
		[]string{
			config.SERVICE_SPEAK + ".DemoService.SayHello",
		})

	log.Printf("sing：%v", res)
	log.Printf("sing：%v", reflect.TypeOf(res))
	log.Printf("sing：%v", res)
	if val := reflect.ValueOf(res); val.IsNil() {
		log.Println("sing：返回值为空")
		return nil
	}

	resp := res.(*proto.DemoResponse)
	if resp.Text == "" {
		log.Println("sing：返回值resp.Text等于空")
		return nil
	}
	rsp.Text = resp.Text

	return nil
}
