package service

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
)

type Option func(opts *Options)

type Mode uint8

const (
	RPC Mode = iota
	HTTP
)

type Options struct {
	Context context.Context

	Name string
	// listen rpc addr
	RpcAddr string
	// listen http addr
	HttpAddr string
	// service registered address
	Advertise string

	Init     []micro.Option
	CallFunc func(micro.Service)

	// use rabbitmq as the broker driver
	Rabbitmq bool

	// set server type ，RPC || HTTP
	ServerType Mode
	Server     server.Server
	BindRoute  func(*gin.Engine)
}

var (
	DEFAULT_RPC_ADDR   = ":8090"
	DEFAULT_HTTP_ADDR  = ":8080"
	DEFAULT_BIND_ROUTE = func(engine *gin.Engine) {
		engine.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "call sxx-micro http server success"})
		})
	}
)

type initKey struct{}
type addressKey struct{}
type callFuncKey struct{}
type rabbitmqKey struct{}
type advertiseKey struct{}
type bindRouteKey struct{}
type serverTypeKey struct{}
type serviceNameKey struct{}

func Name(val string) Option {
	return SetOption(serviceNameKey{}, val)
}

func RabbitmqBroker(val bool) Option {
	return SetOption(rabbitmqKey{}, val)
}

func ServerType(val Mode) Option {
	return SetOption(serverTypeKey{}, val)
}

// when server type is http, bind route for server
func BindRoute(val func(*gin.Engine)) Option {
	return SetOption(bindRouteKey{}, val)
}

// registry node address
func Advertise(val string) Option {
	return SetOption(advertiseKey{}, val)
}

func Init(val []micro.Option) Option {
	return SetOption(initKey{}, val)
}

func CallFunc(val func(service micro.Service)) Option {
	return SetOption(callFuncKey{}, val)
}

func Address(val string) Option {
	return SetOption(addressKey{}, val)
}
