package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"log"
	"os"
	"os/signal"
	proto "sxx-go-micro/proto"
	"sxx-go-micro/Common/wrapper/breaker/hystrix"
	"syscall"
	"time"
)

func main() {
	serviceCirculation()
}

func service() {

	// 客户端从consul中发现服务
	registry := consul.NewRegistry()

	// 创建一个新的服务
	service := micro.NewService(
		// 使用grpc协议
		micro.Client(grpc.NewClient()),
		micro.Name("Greeter.Client"),
		micro.Registry(registry),
	)
	// 初始化
	service.Init()

	// 创建 Greeter 客户端
	greeter := proto.NewGreeterService("Greeter", service.Client())

	// 远程调用 Greeter 服务的 Hello 方法
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "学院君"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}

func serviceCirculation() {
	// 将服务注册到consul
	registry := consul.NewRegistry()

	// hystrix 配置
	hystrix.Configure([]string{"Greeter.test"})

	// 创建一个新的服务
	service := micro.NewService(
		//micro.Transport(grpc.NewTransport()),
		micro.Client(grpc.NewClient()),
		micro.Name("Greeter.Client"),
		micro.Registry(registry),
		// 使用 hystrix
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	// 初始化
	service.Init()

	// 创建 Greeter 客户端
	greeter := proto.NewGreeterService("Greeter", service.Client())

	// 模拟常驻内存
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case os.Interrupt, os.Kill, syscall.SIGQUIT:
				fmt.Println("退出客户端")
				os.Exit(0)
			default:
				fmt.Println("程序执行中...")
			}
		}
	}()

	// 远程调用 Greeter 服务的 Hello 方法
	for {
		rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "学院君"})
		if err != nil {
			log.Fatalf("服务调用失败：%v", err)
			return
		}
		// Print response
		log.Println(rsp.Greeting)
		time.Sleep(3 * time.Second)
	}
}
