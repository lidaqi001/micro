package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/broker/kafka/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/asyncKafka/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	service.Create(config.SERVICE_ASYNC_SUBSCRIBER, func(service micro.Service) {

		pbsb := service.Options().Broker

		if err := pbsb.Connect(); err != nil {
			fmt.Println("broker connection failed!")
			return
		}
		//defer pbsb.Disconnect()

		// 回调
		func(service micro.Service, pbsb broker.Broker) {

			_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.CallSing)

			_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.SingEvent)
		}(service, pbsb)

	}, micro.Broker(
		//kafka.NewBroker(broker.Addrs("kafka-manager.liqi-test.shanxiuxia.com:9092")),
		kafka.NewBroker(broker.Addrs("47.114.98.80:30087")),
	))
}
