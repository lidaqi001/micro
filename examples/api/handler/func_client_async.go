package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"github.com/lidaqi001/micro/examples/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
)

func (h *handler) ClientAsyncA() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			rsp interface{}
			err error
		)
		params := client.Params{
			ClientName: "ginClientAsyncA",
			CallUserFunc: func(srv micro.Service, ctx context.Context, i2 interface{}) (i interface{}, err error) {

				// 业务代码处理
				cli := user.NewDemoService(config.SERVICE_ASYNC_EVENT, srv.Client())
				return cli.SayHelloByUserId(ctx, &user.UserRequest{Id: "ClientAsyncA"})
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

func (h *handler) ClientAsyncB() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			rsp interface{}
			err error
		)
		params := client.Params{
			ClientName: "ginClientAsyncB",
			CallUserFunc: func(srv micro.Service, ctx context.Context, i2 interface{}) (i interface{}, err error) {

				// 业务代码处理
				cli := user.NewDemoService(config.SERVICE_ASYNC_EVENT, srv.Client())
				return cli.SayHello(ctx, &user.DemoRequest{Name: "ClientAsyncB"})
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
