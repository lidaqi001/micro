package main

import (
	"context"
	"encoding/json"
	traceplugin "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/opentracing/opentracing-go"
	"log"
	"github.com/lidaqi001/micro/examples/proto/event"
	"github.com/lidaqi001/micro/plugins/wrapper/service/trace"
	"github.com/lidaqi001/micro/plugins/wrapper/trace/jaeger"
	"time"
)

func main2() {

	service := micro.NewService(
		micro.Name("weather"),
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	pbsb := service.Server().Options().Broker
	pbsb.Connect()
	go func() {
		for now := range time.Tick(5 * time.Second) {
			log.Println("Publishering weather alerts to topic: alerts")
			test := &event.DemoEvent{
				City:        "beijing",
				Timestamp:   now.UTC().Unix(),
				Temperature: 28,
			}
			body, _ := json.Marshal(test)
			pbsb.Publish("alerts", &broker.Message{
				Header: map[string]string{
					"city": test.City,
				},
				Body: body,
			})
		}
	}()
	//p := micro.NewEvent("alerts", service.Client())
	//go func() {
	//	for now := range time.Tick(5 * time.Second) {
	//		log.Println("Publishering weather alerts to topic: alerts")
	//		p.Publish(context.TODO(), &event.DemoEvent{
	//			City:        "beijing",
	//			Timestamp:   now.UTC().Unix(),
	//			Temperature: 28,
	//		})
	//	}
	//}()
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	log.Println("service")
	serviceName := "weather"
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
		//micro.WrapSubscriber(
		//	traceplugin.NewSubscriberWrapper(opentracing.GlobalTracer()),
		//	trace.SubWrapper,
		//),
	)
	service.Init()

	p := micro.NewEvent("alerts", service.Client())
	go func() {
		for now := range time.Tick(3 * time.Second) {
			log.Println("Publishering weather alerts to topic: alerts")
			//test := &event.DemoEvent{
			//	City:        "beijing",
			//	Timestamp:   now.UTC().Unix(),
			//	Temperature: 28,
			//}
			//body, _ := json.Marshal(test)
			//_ = p.Publish(context.TODO(), &broker.Message{
			//	Header: map[string]string{
			//		"city": test.City,
			//		"name": "liqi",
			//	},
			//Body: body,
			//})
			_ = p.Publish(context.TODO(), &event.DemoEvent{
				City:        "beijing",
				Timestamp:   now.UTC().Unix(),
				Temperature: 28,
			})
		}
	}()
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
