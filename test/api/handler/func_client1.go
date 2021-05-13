package handler

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"sxx-go-micro/common/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/plugins/client"
)

func (h handler) Client1() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			rsp interface{}
			err error
			srv micro.Service
			ctx context.Context
		)
		service, ok := c.Get("gin")
		log.Printf("%v", reflect.TypeOf(service))
		if !ok {
			params := client.Params{
				ClientName: "gin",
				HystrixService: []string{
					config.SERVICE_SPEAK + ".DemoService.SayHello",
				},
			}
			client.Create(params, c)
		}
		cc, _ := c.Get("gin_ctx")
		service, ok = c.Get("gin")
		log.Printf("%v", reflect.TypeOf(service))
		srv = service.(micro.Service)
		ctx = cc.(context.Context)

		// 业务代码处理
		cli := user.NewDemoService(config.SERVICE_SING, srv.Client())
		rsp, err = cli.SayHello(ctx, &user.DemoRequest{Name: "李琪"})

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

//func (h handler) Client1() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		input := make(map[string]string)
//		input["a"] = "a"
//		input["b"] = "b"
//		input["c"] = "c"
//
//		params := client.Params{
//			ClientName: "gin",
//			Input:      input,
//			HystrixService: []string{
//				config.SERVICE_SPEAK + ".DemoService.SayHello",
//			},
//			CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
//				// 业务代码处理
//				//i := input.(map[string]string)
//				//log.Printf("传参:::%v,%v", input, i["a"])
//
//				cli := user.NewDemoService(config.SERVICE_SING, svc.Client())
//				req := &user.DemoRequest{Name: "李琪"}
//				return cli.SayHello(ctx, req)
//			},
//		}
//		rsp, err := client.Create(params)
//		code := 200
//		text := ""
//		if err != nil {
//			code = 500
//			text = err.Error()
//		}
//		if text == "" {
//			text = rsp.(*user.DemoResponse).Text
//		}
//		c.JSON(code, gin.H{"message": text})
//	}
//}
