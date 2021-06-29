package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/plugins/api/jwt"
	"github.com/lidaqi001/micro/plugins/api/middleware"
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

	a := &api{opts: options}

	if err := a.init(opts...); err != nil {
		return err
	}

	return a.run()
}

func (a *api) Init(opts ...Option) error {

	if err := a.init(opts...); err != nil {
		return err
	}
	return nil
}

func (a *api) init(opts ...Option) error {

	for _, o := range opts {
		o(&a.opts)
	}

	if val, ok := a.opts.Context.Value(routeKey{}).(func(engine *gin.Engine)); ok {
		a.opts.Router = val
	}

	return nil
}

func (a *api) run() error {

	g := gin.Default()

	// 注册默认路由
	defaultRoute(g)

	// 注册客户端路由
	if a.opts.Router != nil {
		a.opts.Router(g)
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	if err := g.Run(); err != nil {
		log.Println("Api startup failed: ", err)
		return err
	}

	return nil
}

func defaultRoute(g *gin.Engine) {

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
}
