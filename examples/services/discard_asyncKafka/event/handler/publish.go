package handler

import (
	"encoding/json"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"log"
)

func (s *DemoServiceHandler) publishSayHelloByUserId(req *user.UserRequest) error {
	body, _ := json.Marshal(req)
	if err := s.Pbsb.Publish(config.ROCKETMQ_TOPIC_DEFAULT, &broker.Message{
		Header: map[string]string{
			"id": req.Id,
		},
		Body: body,
	}); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}

func (s *DemoServiceHandler) publishSayHello(req string) error {

	body, _ := json.Marshal(req)
	if err := s.Pbsb.Publish(config.ROCKETMQ_TOPIC_DEFAULT, &broker.Message{
		Header: map[string]string{
			"id": req,
		},
		Body: body,
	}); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}
