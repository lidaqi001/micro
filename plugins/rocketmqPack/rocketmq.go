package rocketmqPack

import (
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/broker/rocketmq"
	"github.com/lidaqi001/micro/plugins/service"
)

func Create(serviceName string, registerService func(service micro.Service, pbsb broker.Broker)) {

	_ = service.Create(
		service.Name(serviceName),
		service.CallFunc(func(service micro.Service) {

			pbsb := service.Options().Broker

			if err := pbsb.Connect(); err != nil {
				fmt.Println("broker connection failed!")
				return
			}
			//defer pbsb.Disconnect()

			// 回调
			registerService(service, pbsb)

		}),
		service.Init([]micro.Option{
			micro.Broker(
				// 设置 rocketmq 作为 broker 驱动
				rocketmq.NewBroker(
					rocketmq.Endpoint(helper.GetConfig("ROCKETMQ_ENDPOINT", config.ROCKETMQ_ENDPOINT)),
					rocketmq.AccessKey(helper.GetConfig("ROCKETMQ_ACCESS_KEY", config.ROCKETMQ_ACCESS_KEY)),
					rocketmq.SecretKey(helper.GetConfig("ROCKETMQ_SECRET_KEY", config.ROCKETMQ_SECRET_KEY)),
					rocketmq.InstanceId(helper.GetConfig("ROCKETMQ_INSTANCE_ID", config.ROCKETMQ_INSTANCE_ID)),
				),
			),
		}),
	)

}
