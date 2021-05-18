package hystrix

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/v3/client"
	"github.com/eapache/go-resiliency/retrier"
	"log"
	"net"
	"net/http"
	"time"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{},
	opts ...client.CallOption) error {

	select {
	case <-ctx.Done():
		// context（上下文）已结束

		log.Println("context already canceled！")
		return errors.New("context already canceled！")


	default:
		// hystrix服务治理

		return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
			// 服务重试

			// 初始化retrier，每隔100ms重试一次，总共重试1次
			// PS::: retrier.ConstantBackoff(1, 100*time.Millisecond)
			// retrier 工作模式和 hystrix 类似，在 Run 方法中将待执行的业务逻辑封装到匿名函数传入即可
			r := retrier.New(retrier.ConstantBackoff(0, 100*time.Millisecond), nil)
			err := r.Run(func() error {
				// 将go-micro客户端的重试次数，设置为0
				return c.Client.Call(ctx, req, rsp, client.WithRetries(0))
			})
			return err

		}, func(err error) error {

			// 服务降级
			// 你可以在这里自定义更复杂的服务降级逻辑作为服务熔断的兜底
			log.Printf("hystrix fallback error: %v", err)
			return err

		})
	}
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}

func Configure(names []string) {
	// Hystrix 有默认的参数配置，这里可以针对某些 API 进行自定义配置
	config := hystrix.CommandConfig{
		Timeout:               2000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	}
	configs := make(map[string]hystrix.CommandConfig)
	for _, name := range names {
		configs[name] = config
	}
	hystrix.Configure(configs)

	// 结合 Hystrix Dashboard 将服务状态信息可视化
	// 不使用可视化工具，可以注释掉
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "88"), hystrixStreamHandler)
}
