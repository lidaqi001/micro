package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/examples/api/handler"
	"github.com/lidaqi001/micro/examples/api/jwt"
	"github.com/lidaqi001/micro/examples/api/middleware"
)

func main() {
	r := gin.Default()

	h := handler.NewHandler()

	r.GET("/", h.WelCome())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//jwt
	r.GET("/login/:username/:password", jwt.Login)
	auth := r.Group("auth").Use(middleware.Auth())
	{
		auth.GET("/verify", jwt.Verify)
		auth.GET("/refresh", jwt.Refresh)
		auth.GET("/sayHello", jwt.SayHello)
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	if err := r.Run(); err != nil {
		fmt.Println("网关启动失败！错误：", err)
	}
}
