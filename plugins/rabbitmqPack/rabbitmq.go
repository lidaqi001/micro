package rabbitmqPack

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/broker/rabbitmq/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/plugins/service"
)

func Create(serviceName string, registerService func(service micro.Service, pbsb broker.Broker)) {

	// 设置rabbitmq地址
	rabbitmq.DefaultRabbitURL = helper.GetConfig("RABBITMQ_ADDR", config.RABBITMQ_ADDR)

	_ = service.Create(
		service.Name(serviceName),
		service.CallFunc(func(service micro.Service) {

			pbsb := service.Options().Broker

			if err := pbsb.Connect(); err != nil {
				fmt.Println("broker connection failed!")
				logger.Fatal("broker connection failed!")
				//return
			}

			// 回调
			registerService(service, pbsb)

		}),
		service.Init([]micro.Option{
			micro.Broker(
				// 设置 rabbitmq 为 broker 驱动
				rabbitmq.NewBroker(
					// 设置：Exchange 为持久化
					// If this option is not set, the exchange will be deleted when rabbitmq restarts
					rabbitmq.DurableExchange(),
					// 设置：订阅时创建持久化队列
					//rabbitmq.PrefetchGlobal(),
				),
			),
			micro.WrapCall(func(callFunc client.CallFunc) client.CallFunc {
				return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
					fmt.Println("111")
					return nil
				}
			}),
			micro.WrapSubscriber(func(subscriberFunc server.SubscriberFunc) server.SubscriberFunc {
				return func(ctx context.Context, msg server.Message) error {
					fmt.Println("222")
					return nil
				}
			}),
		}),
	)
}
