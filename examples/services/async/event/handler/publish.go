package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"github.com/lidaqi001/micro/examples/proto/event"
	"github.com/lidaqi001/micro/examples/proto/user"
	"time"
)

func (s *DemoServiceHandler) publishSayHelloByUserId(req *user.UserRequest) error {

	p := micro.NewEvent("singEvent", s.Service.Client())
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

	p := micro.NewEvent("callSing", s.Service.Client())
	if err := p.Publish(context.TODO(), &event.DemoEvent{
		City:        req,
		Timestamp:   time.Now().UTC().Unix(),
		Temperature: 28,
	}); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}
