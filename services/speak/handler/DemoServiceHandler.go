package handler

import (
	"context"
	"log"
	"sxx-go-micro/Common/config"
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
	// 记录响应
	// 在函数返回 stop span 之前，统计函数执行时间
	//defer jaeger.GetTraceServiceSpan(&ctx, jaeger.ServiceSpan{Req: req, Rsp: rsp})

	// 随机休眠时间，模仿实际情况中的慢请求
	//num := rand.Intn(3)
	//time.Sleep(time.Duration(num) * time.Second)

	rsp.Text = config.SERVICE_SPEAK + "::你好, " + req.Name
	log.Println(rsp.Text)
	return nil
}
