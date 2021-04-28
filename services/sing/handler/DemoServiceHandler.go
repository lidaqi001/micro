package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"sxx-go-micro/Common/client"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/Common/helper"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
	"sxx-go-micro/proto"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *proto.UserRequest, *proto.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	sp := jaeger.GetTraceServiceSpan(&ctx, req, nil, nil)

	// 调用 speak
	res, err := client.CreateClient(
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

	// 调用 listen 服务
	res2, err := client.CreateClient(
		"client.2",
		func(service micro.Service, ctx context.Context) (interface{}, interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_LISTEN, service.Client())
			req := &proto.DemoRequest{Name: "学院君2"}
			resp, err := cli.SayHello(ctx, req)
			return req, resp, err
		}, ctx, sp,
		[]string{
			config.SERVICE_LISTEN + ".DemoService.SayHello",
		})

	log.Printf("sing：%v", res)
	log.Printf("listen：%v", res2)

	text := ""
	if !helper.IsNil(res) {
		text += res.(*proto.DemoResponse).Text
	}
	if !helper.IsNil(res2) {
		text += res2.(*proto.DemoResponse).Text
	}
	log.Println("拼接结果：" + text)

	rsp.Text = text

	// 记录响应
	// 在函数返回 stop span 之前，统计函数执行时间
	defer jaeger.SpanSetResService(sp, rsp, err)

	return nil
}
