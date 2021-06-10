package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/proto/user"
	"log"
)

type DemoServiceHandler struct {
	Service micro.Service
	Pbsb    broker.Broker
}

func (s *DemoServiceHandler) SayHelloByUserId(ctx context.Context, req *user.UserRequest, rsp *user.DemoResponse) error {
	text := "async.event::你好, " + req.Id
	log.Println("拼接结果：" + text)
	rsp.Text = text
	go func() {
		_ = s.publishSayHelloByUserId(req)
	}()
	return nil
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {
	text := "async.event::你好, " + req.Name
	log.Println("拼接结果：" + text)
	rsp.Text = text
	go func() {
		_ = s.publishSayHello(text)
	}()
	return nil
}
