package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/services/async/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_ = service.Create(
		service.Name(config.SERVICE_ASYNC_SUBSCRIBER),
		service.CallFunc(func(service micro.Service) {
			_ = micro.RegisterSubscriber("singEvent", service.Server(), handler.SingEvent)
			_ = micro.RegisterSubscriber("singEvent", service.Server(), handler.CallSing)
		}),
	)
}
