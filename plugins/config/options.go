package config

import "context"

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
	DEFAULT_CONFIG_ETCD_ENDPOINT = "http://127.0.0.1:2379"
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
