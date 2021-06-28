package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/plugins/api/jwt"
	"github.com/lidaqi001/micro/plugins/api/middleware"
	config2 "github.com/lidaqi001/micro/plugins/config"
	_ "github.com/spf13/viper/remote"
	"log"
)

type api struct {
	opts Options
}

func Create(opts ...Option) error {
	options := Options{
		Context: context.Background(),
	}

	c := &api{opts: options}

	if err := c.Init(opts...); err != nil {
		return err
	}

	return c.run()
}
func (a *api) Init(opts ...Option) error {

	for _, o := range opts {
		o(&a.opts)
	}

	if val, ok := a.opts.Context.Value(routeKey{}).(func(engine *gin.Engine)); ok {
		a.opts.Router = val
	}
	// set viper remote configuration (https://github.com/spf13/viper)
	if val, ok := a.opts.Context.Value(configPathKey{}).(string); ok && len(val) > 0 {
		a.opts.ConfigPath = val
	}
	if val, ok := a.opts.Context.Value(configTypeKey{}).(string); ok && len(val) > 0 {
		a.opts.ConfigType = val
	}
	if val, ok := a.opts.Context.Value(configEtcdEndpointKey{}).(string); ok && len(val) > 0 {
		a.opts.ConfigEtcdEndpoint = val
	}

	// set viper remote configuration
	if err := config2.RemoteConfig(
		config2.ConfigPath(a.opts.ConfigPath),
		config2.ConfigType(a.opts.ConfigType),
		config2.ConfigEtcdEndpoint(a.opts.ConfigEtcdEndpoint),
	); err != nil {
		return err
	}

	return nil
}

func (a *api) run() error {

	g := gin.Default()

	// welcome
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is SXX Micro Gateway!",
			"version": "Powered By Gin~",
		})
	})

	// jwt（暂时不用，身份验证待完善）
	g.GET("/login/:username/:password", jwt.Login)
	auth := g.Group("auth").Use(middleware.Auth())
	{
		auth.GET("/verify", jwt.Verify)
		auth.GET("/refresh", jwt.Refresh)
		auth.GET("/sayHello", jwt.SayHello)
	}

	// 注册客户端路由
	a.opts.Router(g)

	// 监听并在 0.0.0.0:8080 上启动服务
	if err := g.Run(); err != nil {
		log.Println("Api startup failed: ", err)
		return err
	}

	return nil
}
