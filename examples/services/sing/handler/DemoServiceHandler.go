package handler

import (
	"context"
	"log"
	"sxx-go-micro/examples/proto/user"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *user.UserRequest, *user.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	//sp := jaeger.GetTraceServiceSpan(ctx)
	//
	//// 调用 speak 服务
	//p1 := client.Params{
	//	ClientName: "client.2",
	//	HystrixService: []string{
	//		config.SERVICE_SPEAK + ".DemoService.SayHello",
	//	},
	//	CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
	//		cli := user.NewDemoService(config.SERVICE_SPEAK, svc.Client())
	//		return cli.SayHello(ctx, req)
	//	},
	//	Ctx: ctx,
	//	Sp:  sp,
	//}
	//res, err := client.Create(p1)
	//log.Printf("speak：%v", res)
	//
	//// 调用 listen 服务
	//p2 := client.Params{
	//	ClientName: "client.2",
	//	HystrixService: []string{
	//		config.SERVICE_LISTEN + ".DemoService.SayHello",
	//	},
	//	CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
	//		cli := user.NewDemoService(config.SERVICE_LISTEN, svc.Client())
	//		return cli.SayHello(ctx, req)
	//	},
	//	Ctx: ctx,
	//	Sp:  sp,
	//}
	//res2, err := client.Create(p2)
	//log.Printf("listen：%v", res2)
	//
	text := ""
	//if !helper.IsNil(res) {
	//	text += res.(*user.DemoResponse).Text
	//}
	//if !helper.IsNil(res2) {
	//	text += res2.(*user.DemoResponse).Text
	//}
	text = "sing::你好, " + req.Name
	log.Println("拼接结果：" + text)
	rsp.Text = text

	return nil
	//return err
}
