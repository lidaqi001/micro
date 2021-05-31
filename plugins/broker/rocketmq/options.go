package rocketmq

import (
	"github.com/asim/go-micro/v3/broker"
)

type optionsKeyType struct{}

type subscribeOptionsKeyType struct{}

var (
	DefaultKeyValue = ""

	optionsKey          = optionsKeyType{}
	subscribeOptionsKey = subscribeOptionsKeyType{}
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
	// 您在控制台创建的 Consumer ID(Group ID)
	groupId string
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

func GroupId(id string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.groupId = id
	}
}
