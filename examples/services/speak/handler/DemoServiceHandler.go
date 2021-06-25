package handler

import (
	"context"
	"fmt"
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/examples/proto/user"
	"github.com/spf13/viper"
)

type DemoServiceHandler struct{}

func (s *DemoServiceHandler) SayHelloByUserId(context.Context, *user.UserRequest, *user.DemoResponse) error {
	panic("implement me")
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *user.DemoRequest, rsp *user.DemoResponse) error {

	// 随机休眠时间，模仿实际情况中的慢请求
	//num := rand.Intn(3)
	//time.Sleep(time.Duration(num) * time.Second)
	fmt.Println("speak:viper:",viper.Get("common.database.host"))
	rsp.Text = config.SERVICE_SPEAK + "::你好, " + req.Name
	return nil
}
