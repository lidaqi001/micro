package api

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Option func(opts *Options)

type Options struct {
	Context context.Context

	Router func(engine *gin.Engine)

	// viper configuration (https://github.com/spf13/viper)
	ConfigPath         string
	ConfigType         string
	ConfigEtcdEndpoint string
}

type routeKey struct{}
type configPathKey struct{}
type configTypeKey struct{}
type configEtcdEndpointKey struct{}

func Route(val func(engine *gin.Engine)) Option {
	return SetOption(routeKey{}, val)
}

func ConfigPath(val string) Option {
	return SetOption(configPathKey{}, val)
}

func ConfigType(val string) Option {
	return SetOption(configTypeKey{}, val)
}

func ConfigEtcdEndpoint(val string) Option {
	return SetOption(configEtcdEndpointKey{}, val)
}
