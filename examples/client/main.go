package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"reflect"
	"sxx-go-micro/examples/config"
	"sxx-go-micro/examples/proto/user"
	"sxx-go-micro/plugins/client"
)

func main() {
	//input := []string{"a", "b", "c"}
	input := make(map[string]string)
	input["a"] = "a"
	input["b"] = "b"
	input["c"] = "c"

	params := client.Params{
		ClientName: "client.1",
		Input:      input,
		CallUserFunc: func(svc micro.Service, ctx context.Context, input interface{}) (interface{}, error) {
			// 业务代码处理
			//i := input.(map[string]string)
			//log.Printf("传参:::%v,%v", input, i["a"])

			cli := user.NewDemoService(config.SERVICE_SING, svc.Client())
			req := &user.DemoRequest{Name: "李琪"}
			return cli.SayHello(ctx, req)
		},
	}
	rsp, _ := client.Create(params)

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
