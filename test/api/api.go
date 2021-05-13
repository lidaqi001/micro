package main

import (
	"github.com/gin-gonic/gin"
	"sxx-go-micro/test/api/handler"
	"sxx-go-micro/test/api/jwt"
	"sxx-go-micro/test/api/middleware"
)

func main() {
	r := gin.Default()

	h := handler.NewHandler()

	//jwt
	r.GET("/login/:username/:password", jwt.Login)
	auth := r.Group("auth").Use(middleware.Auth())
	{
		auth.GET("/verify", jwt.Verify)
		auth.GET("/refresh", jwt.Refresh)
		auth.GET("/sayHello", jwt.SayHello)
	}

	r.GET("/", h.WelCome())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/client1", h.Client1())
	call := r.Group("/call")
	{
		call.GET("/client1", h.Client1())
	}

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
