package handler

import (
	"encoding/json"
	"fmt"
	"github.com/asim/go-micro/v3/broker"
	"log"
)

// 因为异步消息没有response对象
// 所以手动记录tarce
//sp := jaeger.NewSpan(ctx)
//sp.SetResponse(out, nil, false)

func CallSing(event broker.Event) error {
	out := fmt.Sprintf("Got sub1：%v", event)
	log.Println(out)
	return nil
}

func SingEvent(event broker.Event) error {
	var rst interface{}
	_ = json.Unmarshal(event.Message().Body, &rst)
	out := fmt.Sprintf("Got sub2：%v", rst)
	log.Println(out)
	return nil
}
