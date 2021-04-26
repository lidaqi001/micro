package main

import (
	"context"
	"github.com/asim/go-micro/v3"
	"log"
	"reflect"
	"sxx-go-micro/Common/client"
	"sxx-go-micro/Common/config"
	"sxx-go-micro/proto"
)

func main() {

	rsp, _ := client.CreateClient(
		"client.1",
		func(service micro.Service, ctx context.Context) (interface{}, interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_SING, service.Client())
			req := &proto.DemoRequest{Name: "学院君"}
			resp, err := cli.SayHello(ctx, req)
			return req, resp, err
		}, nil, nil, )

	log.Printf("%v", reflect.TypeOf(rsp))
	log.Printf("%v", rsp)
	if val := reflect.ValueOf(rsp); val.IsNil() {
		log.Println("返回值为空")
		return
	}

	resp := rsp.(*proto.DemoResponse)
	if resp.Text == "" {
		log.Println("返回值resp.Text等于空")
		return
	}

	log.Println("返回值：" + resp.Text)

}
