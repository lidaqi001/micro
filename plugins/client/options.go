package client

import (
	"context"
	"github.com/asim/go-micro/v3"
)

type Option func(opts *Options)

//type Params struct {
//	ClientName     string
//	HystrixService []string
//	CallUserFunc   func(micro.Service, context.Context, interface{}) (interface{}, error)
//	Ctx            context.Context
//	Input          interface{}
//}
type Options struct {
	// use for client context set option
	Context context.Context

	Ctx      context.Context
	Name     string
	Hystrix  []string
	Init     []micro.Option
	Input    interface{}
	CallFunc func(micro.Service, context.Context, interface{}) (interface{}, error)
}

type ctxKey struct{}
type initKey struct{}
type nameKey struct{}
type inputKey struct{}
type hystrixKey struct{}
type callFuncKey struct{}

func Name(name string) Option {
	return SetOption(nameKey{}, name)
}

func Ctx(c context.Context) Option {
	return SetOption(ctxKey{}, c)
}

func Hystrix(h []string) Option {
	return SetOption(hystrixKey{}, h)
}

func Input(i interface{}) Option {
	return SetOption(inputKey{}, i)
}

func Init(options []micro.Option) Option {
	return SetOption(initKey{}, options)
}

func CallFunc(fn func(micro.Service, context.Context, interface{}) (interface{}, error)) Option {
	return SetOption(callFuncKey{}, fn)
}
