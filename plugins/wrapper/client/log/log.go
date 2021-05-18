package log

import (
	"context"
	"github.com/asim/go-micro/v3/client"
)

// Implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

// log wrapper logs every time a request is made
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	//log.Printf("[wrapper] client request service: %s method: %s\n", req.Service(), req.Endpoint())

	// 请求服务
	err := l.Client.Call(ctx, req, rsp)

	//log.Printf("[wrapper] client rsp: %v\n", rsp)
	return err
}