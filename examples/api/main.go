package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/examples/api/handler"
	"github.com/lidaqi001/micro/plugins/api"
)

func main() {

	//v, _ := config.LoadConfigFromEtcd(config.EtcdEndpoint("http://192.168.1.146:2379"))

	_ = api.Create(api.Route(func(g *gin.Engine) {

		h := handler.NewHandler()

		// 示例请求
		call := g.Group("/call")
		{
			// 同步sing、speak、listen服务相互调用
			call.GET("/handler", h.Client())
			// 异步消息接口A
			call.GET("/asyncA", h.ClientAsyncA())
			// 异步消息接口B
			call.GET("/asyncB", h.ClientAsyncB())
		}
	}))

}
