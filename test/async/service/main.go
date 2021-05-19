package main

import (
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"log"
	"math/rand"
	"sxx-go-micro/examples/proto/event"
	"time"
)

func main() {
	log.Println("service")
	service := micro.NewService(micro.Name("weather"))
	p := micro.NewEvent("alerts", service.Client())
	//go func() {
	for now := range time.Tick(time.Duration(rand.Intn(3)) * time.Second) {
		log.Println("Publishering weather alerts to topic: alerts")
		test := &event.DemoEvent{
			City:        "beijing",
			Timestamp:   now.UTC().Unix(),
			Temperature: 28,
		}
		body, _ := json.Marshal(test)
		_ = p.Publish(context.TODO(), &broker.Message{
			Header: map[string]string{
				"city": test.City,
			},
			Body: body,
		})
	}
	//}()
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
