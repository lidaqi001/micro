package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
	"log"
	"reflect"
)

func main() {
	//input := []string{"a", "b", "c"}
	input := make(map[string]string)
	input["a"] = "a"
	input["b"] = "b"
	input["c"] = "c"

	rsp, _ := client.Create(
		client.Name("client.1"),
		client.Ctx(context.Background()),
		client.Input(input),
		client.CallFunc(func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
			// 业务代码处理
			//i := input.(map[string]string)
			//log.Printf("传参:::%v,%v", input, i["a"])

			cli := user.NewDemoService(config.SERVICE_SING, svc.Client())
			req := &user.DemoRequest{Name: "lidaqi"}
			return cli.SayHello(ctx, req)
		}),
	)

	switch {
	case reflect.ValueOf(rsp).IsNil():
		log.Println("返回值为空")
		return
		//fallthrough
	case rsp.(*user.DemoResponse).Text == "":
		log.Println("返回值resp.Text等于空")
		return
	}

	log.Printf("%v", rsp)
}
