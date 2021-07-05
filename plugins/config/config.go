package config

import (
	"context"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type config struct {
	opts Options
}

const (
	ERROR_DEBUG_ENV = "Check environment variables: DEBUG"

	SUFFIX_DEV  = "_dev"
	SUFFIX_PROD = "_prod"
)

//RemoteConfig
func LoadConfigFromEtcd(opts ...Option) (*viper.Viper, error) {
	options := Options{
		Context: context.Background(),

		ConfigPath:         DEFAULT_CONFIG_PATH,
		ConfigType:         DEFAULT_CONFIG_TYPE,
		ConfigEtcdEndpoint: "http://" + helper.GetRegistryAddress(),
	}

	c := &config{opts: options}

	c.init(opts...)

	return c.run()
}

func (c *config) init(opts ...Option) {

	for _, o := range opts {
		o(&c.opts)
	}

	// set viper remote configuration (https://github.com/spf13/viper)
	if val, ok := c.opts.Context.Value(configPathKey{}).(string); ok && len(val) > 0 {
		c.opts.ConfigPath = val
	}
	if val, ok := c.opts.Context.Value(configTypeKey{}).(string); ok && len(val) > 0 {
		c.opts.ConfigType = val
	}
	if val, ok := c.opts.Context.Value(configEtcdEndpointKey{}).(string); ok && len(val) > 0 {
		c.opts.ConfigEtcdEndpoint = val
	}

}

func (c *config) run() (*viper.Viper, error) {

	var (
		err    error
		suffix string
	)

	// The default configuration is development environment ： {c.opts.ConfigPath}_dev
	suffix = SUFFIX_DEV

	v, err := helper.IsOpenDebug()
	if err == nil {
		if !v {
			// Set to the production environment ： {c.opts.ConfigPath}_prod
			suffix = SUFFIX_PROD
		}
	} else {
		return nil, err
	}

	// alternatively, you can create a new viper instance.
	var runtime_viper = viper.New()

	if err = runtime_viper.AddRemoteProvider(
		"etcd",
		c.opts.ConfigEtcdEndpoint,
		(c.opts.ConfigPath + suffix),
	); err != nil {
		return nil, err
	}

	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	// because there is no file extension in a stream of bytes,
	runtime_viper.SetConfigType(c.opts.ConfigType)

	// read from remote config the first time.
	if err = runtime_viper.ReadRemoteConfig(); err != nil {
		return nil, err
	}

	// listen for remote configuration changes
	if err = runtime_viper.WatchRemoteConfigOnChannel(); err != nil {
		return nil, err
	}

	return runtime_viper, nil
}
