package trace

import (
	"context"
	"github.com/asim/go-micro/v3/server"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

// spanWrapper is a handler wrapper
func SpanWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		//log.Printf("[wrapper] server request: %v", req.Endpoint())
		//log.Printf("[wrapper] server request params: %v", req.Body())

		sp := jaeger.NewSpan(ctx)

		// Trace：请求service前打印请求tag
		sp.SetReq(req.Body())

		// 执行 service 注册函数
		err := fn(ctx, req, rsp)

		// Trace：执行函数后打印 返回值/错误 tag
		sp.SetRes(rsp, err)

		//log.Printf("[wrapper] server rsp: %v", rsp)
		return err
	}
}
