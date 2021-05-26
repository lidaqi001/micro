package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"sxx-go-micro/common/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/plugins/client"
)

func (h *handler) Client() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			rsp interface{}
			err error
		)
		params := client.Params{
			ClientName: "gin",
			CallUserFunc: func(srv micro.Service, ctx context.Context, i2 interface{}) (i interface{}, err error) {
				// 业务代码处理
				cli := user.NewDemoService(config.SERVICE_SING, srv.Client())
				return cli.SayHello(ctx, &user.DemoRequest{Name: "李琪"})
			},
		}
		rsp, err = client.Create(params)

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
