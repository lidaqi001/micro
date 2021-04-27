package service

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/juju/ratelimit"
	"github.com/opentracing/opentracing-go"
	"log"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

func CreateService(ctx context.Context, serviceName string, registerService func(service micro.Service)) {

	// 将服务注册到consul
	registry := consul.NewRegistry()
	// 将服务注册到etcd
	//registry := etcd.NewRegistry()

	// 限流
	bucket := ratelimit.NewBucketWithRate(float64(config.QPS), int64(config.QPS))

	// 初始化全局服务追踪
	// 设置trace server地址
	/*
		1/ 从外部设置env		PS：traceIp := os.Getenv("MICRO_TRACE_SERVER")
		2/ 代码中设置env		PS：traceIp := os.Setenv("MICRO_TRACE_IP", "192.168.1.146")
		参数传递：
			t, io, err := trace.NewTracer("service.trace", traceServer, traceIp)
	*/
	log.Printf("opentracing.GlobalTracer()1:::%v", opentracing.GlobalTracer())
	_, _, err := jaeger.NewTracer(serviceName, config.TRACE_PORT, "")
	if err != nil {
		log.Fatal(err)
	}
	//defer io.Close()
	log.Printf("opentracing.GlobalTracer()2:::%v", opentracing.GlobalTracer())
	//opentracing.SetGlobalTracer(t)
	//log.Printf("opentracing.GlobalTracer()3:::%v", opentracing.GlobalTracer())

	// 创建新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Server(grpc.NewServer()),
		micro.Name(serviceName),
		micro.Registry(registry),

		// 基于ratelimit 限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(bucket, false)),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(
			traceplugin.NewHandlerWrapper(opentracing.GlobalTracer()),
			//traceplugin.NewHandlerWrapper(t),
		),
	)

	// 初始化，会解析命令行参数
	service.Init()

	// 注册处理器，调用服务接口处理请求
	registerService(service)

	// 启动服务
	if err := service.Run(); err != nil {
		log.Println(err)
		fmt.Println(err)
	}
}
