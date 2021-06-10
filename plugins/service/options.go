package service

import (
	"context"
	"github.com/asim/go-micro/v3"
)

type Option func(opts *Options)

type Options struct {
	Context context.Context

	Name     string
	Init     []micro.Option
	CallFunc func(service micro.Service)
}

type initKey struct{}
type callFuncKey struct{}
type serviceNameKey struct{}

func Name(name string) Option {
	return SetOption(serviceNameKey{}, name)
}

func Init(options []micro.Option) Option {
	return SetOption(initKey{}, options)
}

func CallFunc(fn func(service micro.Service)) Option {
	return SetOption(callFuncKey{}, fn)
}
