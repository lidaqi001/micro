package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/proto/event"
	"log"
)

func CallSing(ctx context.Context, event *event.DemoEvent) error {
	//time.Sleep(time.Second * 1)
	out := fmt.Sprintf("Got event1：%v", event)
	log.Println(out)
	return errors.New("event1")
}
func Ha() broker.Handler {
	return func(event broker.Event) error {
		var res interface{}
		json.Unmarshal(event.Message().Body,&res)
		fmt.Println("got 111::", event.Topic(), res)
		return nil
	}
}
func SingEvent() broker.Handler {
	return func(event broker.Event) error {
		var res interface{}
		json.Unmarshal(event.Message().Body,&res)
		fmt.Println("got 111::", event.Topic(), res)
		return nil
	}
}

//func SingEvent(ctx context.Context, event *event.DemoEvent) error {
//	//time.Sleep(time.Second * 20)
//	out := fmt.Sprintf("Got event2：%v", event)
//	log.Println(out)
//	return errors.New("event2")
//}
