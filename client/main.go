package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"reflect"
	"sxx-go-micro/Common/client"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/proto"
)

func main() {

	rsp, _ := client.CreateClient(
		"client.1",
		func(service micro.Service, ctx context.Context) (interface{}, interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_SING, service.Client())
			req := &proto.DemoRequest{Name: "学院君"}
			resp, err := cli.SayHello(ctx, req)
			return req, resp, err
		}, nil, nil,
		[]string{
			config.SERVICE_SING + ".DemoService.SayHello",
		})

	log.Printf("%v", reflect.TypeOf(rsp))
	log.Printf("%v", rsp)
	if val := reflect.ValueOf(rsp); val.IsNil() {
		log.Println("返回值为空")
		return
	}

	resp := rsp.(*proto.DemoResponse)
	if resp.Text == "" {
		log.Println("返回值resp.Text等于空")
		return
	}

	log.Println("返回值：" + resp.Text)

}

//func main() {
//	clientName := "client.1"
//	// 设置trace server地址
//	t, io, err := trace.NewTracer(clientName, config.TRACE_PORT, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer io.Close()
//
//	// 创建空的上下文, 生成追踪 span
//	sp, ctx := opentracing.StartSpanFromContext(context.Background(), "call Client")
//
//	md, ok := metadata.FromContext(ctx)
//	if !ok {
//		md = make(map[string]string)
//	}
//	defer sp.Finish()
//
//	// 注入 opentracing textmap 到空的上下文用于追踪
//	opentracing.GlobalTracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
//	ctx = opentracing.ContextWithSpan(ctx, sp)
//	ctx = metadata.NewContext(ctx, md)
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
//	cli := proto.NewDemoService(config.SERVICE_SING, service.Client())
//	req := &proto.DemoRequest{Name: "学院君"}
//	resp, err := cli.SayHello(ctx, req)
//
//	// 记录请求
//	sp.SetTag("req", req)
//
//	if err != nil {
//		// 记录错误
//		sp.SetTag("err", err)
//		log2.Warnf("服务调用失败：%v", err)
//	}
//
//	// 记录响应
//	sp.SetTag("resp", resp)
//
//	return
//}
