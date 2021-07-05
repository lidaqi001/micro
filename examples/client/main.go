package main

import (
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
	"log"
)

func main() {
	//os.Setenv("REGISTRY_ADDR", "etcd1.shan-service.svc.cluster.local:2379")

	input := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	// 请求服务端
	rsp, _ := client.Create(
		client.Name("client.1"),
		client.Input(input),
		client.CallFunc(func(p client.CallFuncParams) (interface{}, error) {

			// 业务代码处理
			i := p.Input.(map[string]string)
			log.Printf("input:::%v,%v", input, i)

			cli := user.NewDemoService(config.SERVICE_SPEAK, p.Service.Client())
			req := &user.DemoRequest{Name: "lidaqi"}
			return cli.SayHello(p.Ctx, req)
		}),
	)

	// 对返回值类型断言
	switch rsp.(type) {
	case *user.DemoResponse:
		rsp = rsp.(*user.DemoResponse)
		if rsp == nil {
			log.Println("返回值为空")
			return
		}
		log.Printf("%v", rsp)
	}

}
