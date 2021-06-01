package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/async/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_ASYNC_SUBSCRIBER,
		func(service micro.Service) {
			_ = micro.RegisterSubscriber("singEvent", service.Server(), handler.SingEvent)
			_ = micro.RegisterSubscriber("singEvent", service.Server(), handler.CallSing)
		})
}
