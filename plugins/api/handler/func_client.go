package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
)

func (h *handler) Client() gin.HandlerFunc {
	return func(c *gin.Context) {

		rsp, err := client.Create(
			client.Name("gin"),
			client.CallFunc(func(p client.CallFuncParams) (i interface{}, err error) {
				// 业务代码处理
				cli := user.NewDemoService(config.SERVICE_SING, p.Service.Client())
				return cli.SayHello(p.Ctx, &user.DemoRequest{Name: "lidaqi"})
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
