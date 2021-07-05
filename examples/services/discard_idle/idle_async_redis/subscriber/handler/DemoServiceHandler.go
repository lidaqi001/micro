package handler

import (
	"encoding/json"
	"fmt"
	"github.com/asim/go-micro/v3/broker"
	"log"
)

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
