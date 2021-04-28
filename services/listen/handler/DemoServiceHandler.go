package handler

import (
	"context"
	"log"
	"math/rand"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
	"sxx-go-micro/proto"
	"time"
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
	defer jaeger.GetTraceServiceSpan(&ctx, req, rsp, nil)

	// 随机休眠时间，模仿实际情况中的慢请求
	num := rand.Intn(3)
	time.Sleep(time.Duration(num) * time.Second)

	rsp.Text = config.SERVICE_LISTEN + "::你好, " + req.Name
	log.Println("request")
	log.Println("speak")
	log.Println(rsp.Text)
	return nil
}
