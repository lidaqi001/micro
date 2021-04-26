package handler

import (
	"context"
	"log"
	"math/rand"
	"sxx-go-micro/proto"
	"sxx-go-micro/trace"
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
	//defer trace.GetTraceServiceSpan(ctx, req, rsp)
	sp := trace.GetTraceServiceSpan(ctx, req, rsp)
	// 记录请求
	sp.SetTag("req", req)
	defer func() {
		// 记录响应
		sp.SetTag("resp", rsp)
		// 在函数返回 stop span 之前，统计函数执行时间
		sp.Finish()
	}()

	num := rand.Intn(3)
	time.Sleep(time.Duration(num) * time.Second)

	rsp.Text = "server.2::你好, " + req.Name
	log.Println("request")
	log.Println("server-trace2")
	log.Println(rsp.Text)
	return nil

}
