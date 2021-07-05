package main

import (
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/examples/services/speak/handler"
	"github.com/lidaqi001/micro/plugins/service"
)

func main() {
	var err error
	//os.Setenv("REGISTRY_ADDR", "etcd1.shan-service.svc.cluster.local:2379")
	//v, err = config2.LoadConfigFromEtcd()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("v.db.user",v.Get("db.user"))

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
