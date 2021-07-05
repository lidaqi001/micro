package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/plugins/service"
	"github.com/lidaqi001/micro/plugins/service/http/middleware"
)

func main() {
	var err error
	//_ = os.Setenv("REGISTRY_ADDR", "etcd1.shan-service.svc.cluster.local:2379")
	//v, err := config2.LoadConfigFromEtcd()
	//if err != nil {
	//	fmt.Println(err)
	//}

	err = service.Create(
		service.Name(config.SERVICE_HTTP),
		service.ServerType(service.HTTP),
		service.BindRoute(func(e *gin.Engine) {
			demo := e.Group("demo")
			{
				demo.GET("/", func(c *gin.Context) {
					c.JSON(200, gin.H{"msg": "call sxx-micro http server success"})
				}, middleware.Demo())

				demo.GET("/sayHello", func(c *gin.Context) {
					msg := "hello "
					if name, ok := c.GetQuery("name"); ok {
						msg += name
					} else {
						msg += "world~"
					}
					c.JSON(200, gin.H{"msg": msg})
				})
			}
		}),
	)
	fmt.Println(err)
}
