package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"sxx-go-micro/examples/proto/user"
)

type DemoServiceHandler struct {
	Service micro.Service
}

func (s *DemoServiceHandler) SayHelloByUserId(ctx context.Context, req *user.UserRequest, rsp *user.DemoResponse) error {
	go func() {
		_ = s.publishSayHelloByUserId(req)
	}()
	return nil
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {
	text := ""
	text = "sing::你好, " + req.Name
	log.Println("拼接结果：" + text)
	rsp.Text = text
	go func() {
		_ = s.publishSayHello(text)
	}()
	return nil
}
