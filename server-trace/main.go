package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/metadata"
	ratelimit "github.com/juju/ratelimit"
	"log"
	"os"
	proto "sxx-go-micro/proto"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"sxx-go-micro/trace"
	"github.com/opentracing/opentracing-go"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *proto.UserRequest, *proto.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {
	// 从微服务上下文中获取追踪信息
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	var sp opentracing.Span
	wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	// 创建新的 Span 并将其绑定到微服务上下文
	sp = opentracing.StartSpan("SayHello", opentracing.ChildOf(wireContext))
	// 记录请求
	sp.SetTag("req", req)
	defer func() {
		// 记录响应
		sp.SetTag("res", rsp)
		// 在函数返回 stop span 之前，统计函数执行时间
		sp.Finish()
	}()

	rsp.Text = "你好, " + req.Name
	log.Println("request")
	return nil
}

// 每秒钟QPS
const QPS = 1

func main() {

	// 将服务注册到consul
	registry := consul.NewRegistry()
	// 将服务注册到etcd
	//registry := etcd.NewRegistry()

	// 限流
	bucket := ratelimit.NewBucketWithRate(float64(QPS), int64(QPS))

	// 初始化全局服务追踪
	t, io, err := trace.NewTracer("service.trace", os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 创建新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Server(grpc.NewServer()),
		micro.Name("service.trace"),
		micro.Registry(registry),
		// 基于ratelimit 限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(bucket, false)),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// 初始化，会解析命令行参数
	service.Init()

	// 注册处理器，调用 Greeter 服务接口处理请求
	proto.RegisterDemoServiceHandler(service.Server(), new(DemoServiceHandler))

	// 启动服务
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
