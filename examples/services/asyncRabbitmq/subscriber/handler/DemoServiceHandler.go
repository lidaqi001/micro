package handler

import (
	"context"
	"fmt"
	"github.com/lidaqi001/micro/examples/proto/event"
	"log"
)

func CallSing(ctx context.Context, event *event.DemoEvent) error {
	//time.Sleep(time.Second * 1)
	out := fmt.Sprintf("Got sub：%v", event)
	log.Println(out)
	return nil
}

func SingEvent(ctx context.Context, event *event.DemoEvent) error {
	//time.Sleep(time.Second * 1)
	out := fmt.Sprintf("Got sub：%v", event)
	log.Println(out)
	return nil
}
