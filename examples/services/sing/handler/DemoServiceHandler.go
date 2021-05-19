package handler

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"log"
	"sxx-go-micro/common/helper"
	"sxx-go-micro/examples/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/plugins/client"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

type DemoServiceHandler struct {
	PubSub broker.Broker
}

func (s *DemoServiceHandler) SayHelloByUserId(ctx context.Context, req *user.UserRequest, rsp *user.DemoResponse) error {
	log.Println("in SayHelloByUserId")
	go func() {
		_ = s.publishEvent(req)
	}()
	return nil
}

func (s *DemoServiceHandler) publishEvent(req *user.UserRequest) error {
	log.Println("in publishEvent")
	// JSON 编码
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// 构建 broker 消息
	msg := &broker.Message{
		Header: map[string]string{
			"email": req.Id,
		},
		Body: body,
	}
	// 通过 broker 发布消息到消息系统
	if err := s.PubSub.Publish("singEvent", msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {

	// 从微服务上下文中获取追踪信息
	// 创建新的 Span 并将其绑定到微服务上下文
	// 记录请求
	sp := jaeger.GetTraceServiceSpan(ctx)

	common := client.Params{
		Sp:         sp,
		Ctx:        ctx,
		ClientName: "client.2",
	}

	// 调用 speak 服务
	p1 := common
	p1.CallUserFunc = func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
		cli := user.NewDemoService(config.SERVICE_SPEAK, svc.Client())
		return cli.SayHello(ctx, req)
	}
	res, err := client.Create(p1)
	log.Printf("speak：%v", res)

	// 调用 listen 服务
	p2 := common
	p2.CallUserFunc = func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
		cli := user.NewDemoService(config.SERVICE_LISTEN, svc.Client())
		return cli.SayHello(ctx, req)
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
