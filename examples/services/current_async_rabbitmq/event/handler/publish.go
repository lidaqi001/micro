package handler

import (
	"context"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/event"
	"github.com/lidaqi001/micro/examples/proto/user"
	e "github.com/lidaqi001/micro/plugins/event"
	"log"
	"time"
)

func (s *DemoServiceHandler) publishSayHelloByUserId(req *user.UserRequest) error {

	// 发布消息
	e := e.New{Client: s.Service.Client()}
	err := e.Publish(config.EVENT_A, context.TODO(), &event.DemoEvent{
		City:        req.Id,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	})
	if err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}

func (s *DemoServiceHandler) publishSayHello(req string) error {

	// 发布消息
	e := e.New{Client: s.Service.Client()}
	err := e.Publish(config.EVENT_B, context.TODO(), &event.DemoEvent{
		City:        req,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	})
	if err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}
