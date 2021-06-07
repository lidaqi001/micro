package rabbitmqPack

import (
	"fmt"
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/service"
)

func Create(serviceName string, registerService func(service micro.Service, pbsb broker.Broker)) {

	// 设置rabbitmq地址
	rabbitmq.DefaultRabbitURL = helper.GetConfig("RABBITMQ_ADDR", config.RABBITMQ_ADDR)

	service.Create(serviceName, func(service micro.Service) {

		pbsb := service.Options().Broker

		if err := pbsb.Connect(); err != nil {
			fmt.Println("broker connection failed!")
			logger.Fatal("broker connection failed!")
			//return
		}

		// 回调
		registerService(service, pbsb)

	}, micro.Broker(
		// 设置 rabbitmq 为 broker 驱动
		rabbitmq.NewBroker(
			// 设置：Exchange 为持久化
			rabbitmq.DurableExchange(),
			// 设置：订阅时创建持久化队列
			//rabbitmq.PrefetchGlobal(),

			//rabbitmq.ExchangeName(serviceName),
		),
	))
}
