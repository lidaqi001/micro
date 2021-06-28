package config

import (
	"context"
	"github.com/lidaqi001/micro/common/helper"
)

type Option func(opts *Options)

type Options struct {
	Context context.Context

	// viper configuration (https://github.com/spf13/viper)
	ConfigPath         string
	ConfigType         string
	ConfigEtcdEndpoint string
}

const (
	DEFAULT_CONFIG_PATH          = "/config"
	DEFAULT_CONFIG_TYPE          = "yaml"
)

type configPathKey struct{}
type configTypeKey struct{}
type configEtcdEndpointKey struct{}

func ConfigPath(val string) Option {
	return SetOption(configPathKey{}, val)
}

func ConfigType(val string) Option {
	return SetOption(configTypeKey{}, val)
}

func ConfigEtcdEndpoint(val string) Option {
	return SetOption(configEtcdEndpointKey{}, val)
}
