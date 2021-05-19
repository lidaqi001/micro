package main

import (
	"github.com/asim/go-micro/v3"
	"log"
	"sxx-go-micro/test/async/client/handler"
)

func main() {
	log.Println("client")
	service := micro.NewService(micro.Name("weather_client"))
	service.Init()
	_ = micro.RegisterSubscriber("alerts", service.Server(), handler.ProcessEvent)
	if err := service.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
