package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/lidaqi001/micro/examples/proto/event"
	"log"
)

func CallSing(ctx context.Context, event *event.DemoEvent) error {
	//time.Sleep(time.Second * 1)
	out := fmt.Sprintf("Got event1：%v", event)
	log.Println(out)
	return errors.New("event1")
}

func SingEvent(ctx context.Context, event *event.DemoEvent) error {
	//time.Sleep(time.Second * 20)
	out := fmt.Sprintf("Got event2：%v", event)
	log.Println(out)
	return errors.New("event2")
}
