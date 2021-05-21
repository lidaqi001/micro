package handler

import (
	"encoding/json"
	"github.com/asim/go-micro/v3/broker"
	"log"
	"sxx-go-micro/examples/proto/user"
)

type DemoServiceHandler struct {
	PubSub broker.Broker
}

func ProcessEvent(event broker.Event) error {
	out := new(user.UserRequest)
	_ = json.Unmarshal(event.Message().Body, &out)
	log.Println("Got event::: ",event.Message().Header, out)
	return nil
}
