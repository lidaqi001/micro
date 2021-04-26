package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"reflect"
	"sxx-go-micro/Common/client"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/proto"
	"sxx-go-micro/trace"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *proto.UserRequest, *proto.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	sp := trace.GetTraceServiceSpan(ctx, req, nil)

	// 记录响应
	// 在函数返回 stop span 之前，统计函数执行时间
	//defer trace.SpanSetResponse(sp, rsp)

	// 记录请求
	sp.SetTag("req", req)
	defer func() {
		// 记录响应
		sp.SetTag("resp", rsp)
		// 在函数返回 stop span 之前，统计函数执行时间
		sp.Finish()
	}()

	// 调用service2
	res, _ := client.CreateClient(
		"client.2",
		func(service micro.Service, ctx context.Context) (interface{}, interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_SPEAK, service.Client())
			req := &proto.DemoRequest{Name: "学院君"}
			resp, err := cli.SayHello(ctx, req)
			return req, resp, err
		}, ctx, sp)

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

	//resp := res.(*proto.DemoResponse)
	//rsp.Text = resp.Text

	return nil
}
