package main

import (
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/speak/handler"
	config2 "github.com/lidaqi001/micro/plugins/config"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	_, err := config2.LoadConfigFromEtcd(config2.EtcdEndpoint("http://192.168.1.146:2379"))
	if err != nil {
		fmt.Println(err)
	}
	err = service.Create(
		service.Name(config.SERVICE_SPEAK),
		//service.Advertise("127.0.0.1:9207"),
		service.CallFunc(func(service micro.Service) {
			// 注册处理函数
			_ = user.RegisterDemoServiceHandler(service.Server(), new(handler.DemoServiceHandler))
		}),
	)
	fmt.Println(err)
}
