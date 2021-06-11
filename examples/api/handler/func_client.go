package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
)

func (h *handler) Client() gin.HandlerFunc {
	return func(c *gin.Context) {

		rsp, err := client.Create(
			client.Name("gin"),
			client.CallFunc(func(srv micro.Service, ctx context.Context, i2 interface{}) (i interface{}, err error) {
				// 业务代码处理
				cli := user.NewDemoService(config.SERVICE_SING, srv.Client())
				return cli.SayHello(ctx, &user.DemoRequest{Name: "lidaqi"})
			}),
		)

		code := 200
		text := ""
		if err != nil {
			code = 500
			text = err.Error()
		}
		if text == "" {
			text = rsp.(*user.DemoResponse).Text
		}
		c.JSON(code, gin.H{"message": text})
	}
}
