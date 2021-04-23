package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/opentracing/opentracing-go"
	"log"
	"os"
	"os/signal"
	proto "sxx-go-micro/proto"
	"sxx-go-micro/Common/wrapper/breaker/hystrix"
	"syscall"
	"time"
	"sxx-go-micro/trace"
)

func main() {
	service()
}

func service() {

	// 客户端从consul中发现服务
	registry := consul.NewRegistry()

	// 初始化追踪器
	//t, io, err := trace.NewTracer("laracom.demo.cli", os.Getenv("MICRO_TRACE_SERVER"))
	traceServer := "192.168.1.145:6831"
	t, io, err := trace.NewTracer("cli.trace", traceServer)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()


	// 创建一个新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		micro.Name("Greeter.Client"),
		micro.Registry(registry),
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
	)
	// 初始化
	service.Init()

	client := proto.NewDemoService("service.trace", service.Client())

	// 创建空的上下文, 生成追踪 span
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "call")
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	defer span.Finish()

	// 注入 opentracing textmap 到空的上下文用于追踪
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)
	// 记录请求 && 响应 && 错误
	req := &proto.DemoRequest{Name: "学院君"}
	span.SetTag("req", req)
	resp, err := client.SayHello(ctx, req)
	if err != nil {
		span.SetTag("err", err)
		log.Fatalf("服务调用失败：%v", err)
		return
	}
	span.SetTag("resp", resp)
	log.Println(resp.Text)

}

func serviceCirculation() {
	// 将服务注册到consul
	registry := consul.NewRegistry()

	// hystrix 配置
	hystrix.Configure([]string{"Greeter.test"})

	// 创建一个新的服务
	service := micro.NewService(
		//micro.Transport(grpc.NewTransport()),
		micro.Client(grpc.NewClient()),
		micro.Name("Greeter.Client"),
		micro.Registry(registry),
		// 使用 hystrix
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	// 初始化
	service.Init()

	// 创建 Greeter 客户端
	greeter := proto.NewDemoService("Greeter", service.Client())

	// 模拟常驻内存
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case os.Interrupt, os.Kill, syscall.SIGQUIT:
				fmt.Println("退出客户端")
				os.Exit(0)
			default:
				fmt.Println("程序执行中...")
			}
		}
	}()

	// 远程调用 Greeter 服务的 Hello 方法
	for {
		rsp, err := greeter.SayHello(context.TODO(), &proto.DemoRequest{Name: "学院君"})
		if err != nil {
			log.Fatalf("服务调用失败：%v", err)
			return
		}
		// Print response
		log.Println(rsp.Text)
		time.Sleep(3 * time.Second)
	}
}
