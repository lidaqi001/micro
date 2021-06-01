package rocketmq

import (
	"github.com/asim/go-micro/v3/broker"
)

type optionsKeyType struct{}
type sOptionsKeyType struct{}
type pOptionsKeyType struct{}

var (
	DefaultKeyValue = ""
	DefaultTagValue = "go-micro"

	optionsKey = optionsKeyType{}

	sOptionsKey = sOptionsKeyType{}

	pOptionsKey = pOptionsKeyType{}
)

// options contain additional options for the broker.
type brokerOptions struct {
	// 设置HTTP接入域名（此处以公共云生产环境为例）
	endpoint string
	// AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
	accessKey string
	// SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
	secretKey string
	// Topic所属实例ID，默认实例为空
	instanceId string
}
type sBrokerOptions struct {
	// 您在控制台创建的 Consumer ID(Group ID)
	groupId string
	// 订阅Tag
	subTag string
}
type pBrokerOptions struct {
	// 消息Tag
	tag string
}

func AccessKey(key string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.accessKey = key
	}
}

func Endpoint(e string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.endpoint = e
	}
}

func SecretKey(key string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.secretKey = key
	}
}

func InstanceId(id string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.instanceId = id
	}
}

func GroupId(id string) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		bo := o.Context.Value(sOptionsKey).(*sBrokerOptions)
		bo.groupId = id
	}
}

func SubTag(key string) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		bo := o.Context.Value(sOptionsKey).(*sBrokerOptions)
		bo.subTag = key
	}
}

func Tag(key string) broker.PublishOption {
	return func(o *broker.PublishOptions) {
		bo := o.Context.Value(pOptionsKey).(*pBrokerOptions)
		bo.tag = key
	}
}
