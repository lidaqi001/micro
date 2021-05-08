package handler

import (
	"context"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/proto"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *proto.UserRequest, *proto.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *proto.DemoRequest, rsp *proto.DemoResponse) error {

	// 随机休眠时间，模仿实际情况中的慢请求
	//num := rand.Intn(3)
	//time.Sleep(time.Duration(num) * time.Second)

	rsp.Text = config.SERVICE_SPEAK + "::你好, " + req.Name
	//log.Println(rsp.Text)
	return nil
}
