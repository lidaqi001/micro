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

	params := client.Params{
		ClientName: "client.1",
		CallUserFunc: func(service micro.Service, ctx context.Context) (interface{}, error) {
			cli := proto.NewDemoService(config.SERVICE_SING, service.Client())
			req := &proto.DemoRequest{Name: "李琪"}
			return cli.SayHello(ctx, req)
		},
		HystrixService: []string{
			config.SERVICE_SING + ".DemoService.SayHello",
		}}
	rsp, _ := client.Create(&params)

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

}
