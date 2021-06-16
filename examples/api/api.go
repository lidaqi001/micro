package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/examples/api/handler"
	"github.com/lidaqi001/micro/examples/api/jwt"
	"github.com/lidaqi001/micro/examples/api/middleware"
	"log"
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
		// 同步sing、speak、listen服务相互调用
		call.GET("/handler", h.Client())
		// 异步消息接口A
		call.GET("/asyncA", h.ClientAsyncA())
		// 异步消息接口B
		call.GET("/asyncB", h.ClientAsyncB())
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	if err := r.Run(); err != nil {
		log.Println("Api startup failed: ", err)
	}
}
