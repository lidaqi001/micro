package config

import (
	"context"
	"fmt"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"time"
)

type config struct {
	opts Options
}

func RemoteConfig(opts ...Option) error {
	options := Options{
		Context:            context.Background(),
		ConfigPath:         "/config",
		ConfigEtcdEndpoint: "http://" + helper.GetRegistryAddress(),
	}

	c := &config{opts: options}

	c.Init(opts...)

	return c.run()
}
func (c *config) Init(opts ...Option) {

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

func (c *config) run() error {

	var err error

	if err = viper.AddRemoteProvider(
		"etcd",
		c.opts.ConfigEtcdEndpoint,
		c.opts.ConfigPath,
	); err != nil {
		return err
	}

	viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes,
	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	if err = viper.ReadRemoteConfig(); err != nil {
		return err
	}

	fmt.Println(
		err,
		viper.AllKeys(),
		viper.Get("common.database.host"),
	)

	// 监听远程配置
	go func() {
		for {
			time.Sleep(time.Second * 1) // 每次请求后延迟一下
			if err := viper.WatchRemoteConfig(); err != nil {
				fmt.Printf("viper.WatchRemoteConfig() error:%v", err)
			}
		}
	}()

	fmt.Println(
		viper.Get("common.database.host"),
	)
	//time.Sleep(time.Second * 10)
	fmt.Println(
		viper.Get("common.database.host"),
	)

	return nil
}

func (c *config) run2() error {

	var err error

	// alternatively, you can create a new viper instance.
	var runtime_viper = viper.New()

	err = runtime_viper.AddRemoteProvider(
		"etcd",
		c.opts.ConfigEtcdEndpoint,
		c.opts.ConfigPath,
	)
	if err != nil {
		return err
	}

	runtime_viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes,
	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err = runtime_viper.ReadRemoteConfig()

	fmt.Println(
		err,
		runtime_viper.AllKeys(),
		viper.Get("common.database.host"),
	)
	// 监听远程配置
	_ = runtime_viper.WatchRemoteConfigOnChannel()

	fmt.Println(
		runtime_viper.Get("common.database.host"),
	)
	time.Sleep(time.Second * 10)
	fmt.Println(
		runtime_viper.Get("common.database.host"),
	)
	if err != nil {
		return err
	}

	return nil
}
