package main

import (
	"context"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/proto/event"
	"github.com/lidaqi001/micro/plugins/wrapper/service/trace"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"github.com/opentracing/opentracing-go"
	"log"
)

func main2() {
	serviceName := "weather_client"
	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	// 设置为全局 tracer
	// PS：只在 service 声明
	opentracing.SetGlobalTracer(t)
	service := micro.NewService(
		micro.Name(serviceName),
		micro.WrapSubscriber(trace.SubWrapper),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 链路追踪中间件
		micro.WrapHandler(trace.SpanWrapper),
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(opentracing.GlobalTracer()),
			trace.SubWrapper,
		),
	)
	service.Init()
	srv := service.Server()
	srv.Init()
	//service.Server().Options().Broker.Init(broker.Addrs(":8000"))
	service.Server().Options().Broker.Subscribe("alerts", processEvent2)

	//broker.Connect()
	//defer broker.Disconnect()
	//broker.Subscribe("alerts", processEvent2)

	pbsb := srv.Options().Broker
	pbsb.Connect()
	defer pbsb.Disconnect()
	pbsb.Subscribe("alerts", processEvent2)
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	log.Println("client")
	serviceName := "weather_client"
	// 初始化全局服务追踪
	t, io, err := jaeger.NewTracer(serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	// 设置为全局 tracer
	// PS：只在 service 声明
	opentracing.SetGlobalTracer(t)
	service := micro.NewService(
		micro.Name(serviceName),
		// 基于 jaeger 采集追踪数据
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 链路追踪中间件
		micro.WrapHandler(trace.SpanWrapper),
		micro.WrapSubscriber(
			traceplugin.NewSubscriberWrapper(opentracing.GlobalTracer()),
			trace.SubWrapper,
		),
	)
	service.Init()
	_ = micro.RegisterSubscriber("alerts", service.Server(), processEvent)
	//_ = micro.RegisterSubscriber("alerts", service.Server(), func(ctx context.Context, event *event.DemoEvent) error {
	//	log.Println("Got alert:", event)
	//	return nil
	//})
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
func processEvent(ctx context.Context, event *event.DemoEvent) error {
	log.Println("Got alert:", event)
	sp := jaeger.NewSpan(ctx)
	sp.SetResponse("123456", nil, false)
	return nil
}
func processEvent1(ctx context.Context, event *broker.Message) error {
	log.Println("Got alert:", event)
	return nil
}

func processEvent2(event broker.Event) error {
	log.Println("Got event::: ", event, event.Message().Header)
	return nil
}
