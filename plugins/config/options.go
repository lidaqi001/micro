package config

import (
	"context"
)

type Option func(opts *Options)

type COption struct {
	// viper configuration (https://github.com/spf13/viper)
	ConfigPath         string
	ConfigType         string
	ConfigEtcdEndpoint string
}

type Options struct {
	// viper configuration (https://github.com/spf13/viper)
	ConfigPath         string
	ConfigType         string
	ConfigEtcdEndpoint string
	//COption
	Context context.Context
}

const (
	DEFAULT_CONFIG_PATH = "/config"
	DEFAULT_CONFIG_TYPE = "yaml"
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

func EtcdEndpoint(val string) Option {
	return SetOption(configEtcdEndpointKey{}, val)
}
