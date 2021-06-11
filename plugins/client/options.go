package client

import (
	"context"
	"github.com/asim/go-micro/v3"
)

type Option func(opts *Options)

type Options struct {
	// use for client context set option
	Context context.Context

	// handler context
	Ctx context.Context

	// ps : []string{"sing.DemoService.SayHello"}
	Hystrix []string

	// CallFunc additional parameters
	Input interface{}

	Name     string
	Init     []micro.Option
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
