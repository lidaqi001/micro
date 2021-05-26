package main

import (
	"github.com/asim/go-micro/v3"
	"sxx-go-micro/common/config"
	"sxx-go-micro/examples/services/async/subscriber/process"
	"sxx-go-micro/plugins/service"
)

func main() {
	service.Create(
		config.SERVICE_ASYNC_SUBSCRIBER,
		func(service micro.Service) {
			_ = micro.RegisterSubscriber("singEvent", service.Server(), process.SingEvent)
			_ = micro.RegisterSubscriber("callSing", service.Server(), process.CallSing)
		})
}
