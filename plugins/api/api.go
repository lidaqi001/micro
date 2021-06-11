package main

import (
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

	// jwt
	r.GET("/login/:username/:password", jwt.Login)
	auth := r.Group("auth").Use(middleware.Auth())
	{
		auth.GET("/verify", jwt.Verify)
		auth.GET("/refresh", jwt.Refresh)
		auth.GET("/sayHello", jwt.SayHello)
	}

	// 示例请求
	call := r.Group("/call")
	{
		call.GET("/handler", h.Client())
		call.GET("/asyncA", h.ClientAsyncA())
		call.GET("/asyncB", h.ClientAsyncB())
	}

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
