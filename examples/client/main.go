package main

import (
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/lidaqi001/micro/plugins/client"
	"log"
)

func main() {
	//input := []string{"a", "b", "c"}
	input := make(map[string]string)
	input["a"] = "a"
	input["b"] = "b"
	input["c"] = "c"

	rsp, _ := client.Create(
		client.Name("client.1"),
		client.Input(input),
		client.CallFunc(func(p client.CallFuncParams) (interface{}, error) {
			// 业务代码处理
			//i := p.Input.(map[string]string)
			//log.Printf("传参:::%v,%v", input, i["a"])

			cli := user.NewDemoService(config.SERVICE_SING, p.Service.Client())
			req := &user.DemoRequest{Name: "lidaqi"}
			return cli.SayHello(p.Ctx, req)
		}),
	)

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
