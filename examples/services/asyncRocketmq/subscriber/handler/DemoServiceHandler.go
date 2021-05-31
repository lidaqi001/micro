package handler

import (
	"fmt"
	"github.com/asim/go-micro/v3/broker"
	"log"
)

// 因为异步消息没有response对象
// 所以手动记录tarce
//sp := jaeger.NewSpan(ctx)
//sp.SetResponse(out, nil, false)

func CallSing(event broker.Event) error {
	out := fmt.Sprintf("Got sub：%v", event)
	log.Println(out)
	return nil
}

func SingEvent(event broker.Event) error {
	out := fmt.Sprintf("Got sub：%v", event)
	log.Println(out)
	return nil
}
