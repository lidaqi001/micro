package main

import (
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/broker"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	c "github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/services/asyncRocketmq/subscriber/handler"
	"github.com/lidaqi001/micro/plugins/asyncRocketmq"
	"github.com/lidaqi001/micro/plugins/broker/rocketmq"
)

func main() {
	asyncRocketmq.Create(
		c.SERVICE_ASYNC_SUBSCRIBER_ROCKETMQ,
		func(service micro.Service, pbsb broker.Broker) {

			c2 := []broker.SubscribeOption{
				rocketmq.GroupId(helper.GetConfig("ROCKETMQ_GROUP_ID", config.ROCKETMQ_GROUP_ID)),
				rocketmq.SubTag("b"),
			}
			_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.CallSing, c2...)

			c := []broker.SubscribeOption{
				rocketmq.GroupId(helper.GetConfig("ROCKETMQ_GROUP_ID", config.ROCKETMQ_GROUP_ID)),
				rocketmq.SubTag("a"),
			}
			_, _ = pbsb.Subscribe(config.ROCKETMQ_TOPIC_DEFAULT, handler.SingEvent, c...)
		})
}
