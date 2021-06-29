package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"strconv"
	"time"
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

		// 配置路径（配置文件必须提前配好）
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

func (c *config) run2() error {

	var (
		err    error
		debug  string
		suffix string
	)

	// 默认配置为开发环境 ： {c.opts.ConfigPath}_dev
	suffix = SUFFIX_DEV

	if debug = os.Getenv("DEBUG"); len(debug) > 0 {

		d, err := strconv.ParseInt(debug, 10, 64)
		if err != nil {
			return errors.New(
				fmt.Sprintf("%s\nerror:%v", ERROR_DEBUG_ENV, err),
			)
		}
		if d == 0 {
			// 设置为生产环境 ： {c.opts.ConfigPath}_prod
			suffix = SUFFIX_PROD
		}
	}

	if err = viper.AddRemoteProvider(
		"etcd",
		c.opts.ConfigEtcdEndpoint,
		(c.opts.ConfigPath + suffix),
	); err != nil {
		return err
	}

	viper.SetConfigType(c.opts.ConfigType) // because there is no file extension in a stream of bytes,
	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	if err = viper.ReadRemoteConfig(); err != nil {
		return err
	}

	go func() {
		for {
			// 每次请求后延迟1s
			time.Sleep(time.Second * 1)

			// 监听远程配置变化
			if err := viper.WatchRemoteConfig(); err != nil {
				fmt.Printf("viper.WatchRemoteConfig() error:%v", err)
			}
		}
	}()

	return nil
}

func (c *config) run() (*viper.Viper, error) {

	var (
		err    error
		debug  string
		suffix string
	)

	// 默认配置为开发环境 ： {c.opts.ConfigPath}_dev
	suffix = SUFFIX_DEV

	if debug = os.Getenv("DEBUG"); len(debug) > 0 {

		d, err := strconv.ParseInt(debug, 10, 64)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("%s\nerror:%v", ERROR_DEBUG_ENV, err),
			)
		}
		if d == 0 {
			// 设置为生产环境 ： {c.opts.ConfigPath}_prod
			suffix = SUFFIX_PROD
		}
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

	runtime_viper.SetConfigType(c.opts.ConfigType) // because there is no file extension in a stream of bytes,
	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	if err = runtime_viper.ReadRemoteConfig(); err != nil {
		return nil, err
	}

	// 监听远程配置变更
	if err = runtime_viper.WatchRemoteConfigOnChannel(); err != nil {
		return nil, err
	}

	return runtime_viper, nil
}
