package handler

import (
	"context"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
	"log"
)

type DemoServiceHandler struct {
}

func (s *DemoServiceHandler) SayHelloByUserId(ctx context.Context, req *user.UserRequest, rsp *user.DemoResponse) error {
	log.Println("hello~")
	return nil
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {

	// 调用 speak 服务
	res, err := client.Create(
		client.Ctx(ctx),
		client.Name("client.2"),
		client.CallFunc(func(p client.CallFuncParams) (interface{}, error) {
			cli := user.NewDemoService(config.SERVICE_SPEAK, p.Service.Client())
			return cli.SayHello(p.Ctx, req)
		}),
	)

	log.Printf("speak：%v", res)

	// 调用 listen 服务
	res2, err := client.Create(
		client.Ctx(ctx),
		client.Name("client.2"),
		client.CallFunc(func(p client.CallFuncParams) (interface{}, error) {
			cli := user.NewDemoService(config.SERVICE_LISTEN, p.Service.Client())
			return cli.SayHello(p.Ctx, req)
		}),
	)
	log.Printf("listen：%v", res2)

	text := ""
	if !helper.IsNil(res) {
		text += res.(*user.DemoResponse).Text
	}
	if !helper.IsNil(res2) {
		text += res2.(*user.DemoResponse).Text
	}
	log.Println("拼接结果：" + text)
	rsp.Text = text

	return err
}

//func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {
//	text := ""
//	text = "sing::你好, " + req.Name
//	log.Println("拼接结果：" + text)
//	rsp.Text = text
//	return nil
//}
