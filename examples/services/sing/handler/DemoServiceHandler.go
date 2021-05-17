package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"sxx-go-micro/common/config"
	"sxx-go-micro/common/helper"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/plugins/client"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *user.UserRequest, *user.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	sp := jaeger.GetTraceServiceSpan(ctx)

	// 调用 speak 服务
	p1 := client.Params{
		Ctx: ctx,
		Sp:  sp,
		ClientName: "client.2",
		CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
			cli := user.NewDemoService(config.SERVICE_SPEAK, svc.Client())
			return cli.SayHello(ctx, req)
		},
	}
	res, err := client.Create(p1)
	log.Printf("speak：%v", res)

	// 调用 listen 服务
	p2 := client.Params{
		Ctx: ctx,
		Sp:  sp,
		ClientName: "client.2",
		CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
			cli := user.NewDemoService(config.SERVICE_LISTEN, svc.Client())
			return cli.SayHello(ctx, req)
		},
	}
	res2, err := client.Create(p2)
	log.Printf("listen：%v", res2)

	text := ""
	if !helper.IsNil(res) {
		text += res.(*user.DemoResponse).Text
	}
	if !helper.IsNil(res2) {
		text += res2.(*user.DemoResponse).Text
	}
	log.Println("拼接结果：" + text)
	rsp.Text = text

	return err
}

//func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {
//	text := ""
//	text = "sing::你好, " + req.Name
//	log.Println("拼接结果：" + text)
//	rsp.Text = text
//	return nil
//}
