package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/event"
	"github.com/lidaqi001/micro/examples/proto/user"
	"log"
	"time"
)

func (s *DemoServiceHandler) publishSayHelloByUserId(req *user.UserRequest) error {

	// 发布消息
	p := micro.NewEvent(config.EVENT_A, s.Service.Client())
	if err := p.Publish(context.TODO(), &event.DemoEvent{
		City:        req.Id,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	}); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}

func (s *DemoServiceHandler) publishSayHello(req string) error {

	// 发布消息
	p := micro.NewEvent(config.EVENT_B, s.Service.Client())
	if err := p.Publish(context.TODO(), &event.DemoEvent{
		City:        req,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	}); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}
