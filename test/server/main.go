package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	"github.com/asim/go-micro/v3"
	ratelimit "github.com/juju/ratelimit"
	"log"
	proto "sxx-go-micro/examples/proto"
)

type GreeterServiceHandler struct{}

// 接口实现
func (g *GreeterServiceHandler) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = " 你好, " + req.Name
	log.Println("接口调用")
	return nil
}

// 每秒钟QPS
const QPS = 1

func main() {

	// 将服务注册到consul
	registry := consul.NewRegistry()
	// 将服务注册到etcd
	//registry := etcd.NewRegistry()

	// 限流
	bucket := ratelimit.NewBucketWithRate(float64(QPS), int64(QPS))

	// 创建新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Server(grpc.NewServer()),
		micro.Name("Greeter"),
		micro.Registry(registry),
		// 基于ratelimit 限流
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(bucket, false)),
	)

	// 初始化，会解析命令行参数
	service.Init()

	// 注册处理器，调用 Greeter 服务接口处理请求
	proto.RegisterGreeterHandler(service.Server(), new(GreeterServiceHandler))

	// 启动服务
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
