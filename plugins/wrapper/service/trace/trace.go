package trace

import (
	"context"
	"github.com/asim/go-micro/v3/server"
	"log"
	"sxx-go-micro/plugins/wrapper/trace/jaeger"
)

// SpanWrapper is a handler wrapper
func SpanWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		log.Printf(
			"[Han Wrapper] server request: %v，server request params: %v",
			req.Endpoint(),
			req.Body(),
		)

		sp := jaeger.NewSpan(ctx)

		// 处理前
		// Trace：记录 请求值
		sp.SetRequest(req.Body())

		err := fn(ctx, req, rsp)

		// 处理后
		// Trace：记录 返回值/错误
		sp.SetResponse(rsp, err)

		log.Printf("[Han Wrapper] server rsp: %v", rsp)

		return err
	}
}

// SubWrapper is a subscriber wrapper
func SubWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {

		log.Printf("[Sub Wrapper] Before serving publication topic: %v", msg.Topic())

		sp := jaeger.NewSpan(ctx)
		sp.SetTopic(msg.Topic())
		sp.SetHeader(msg.Header())
		sp.SetPayload(msg.Payload())

		err := fn(ctx, msg)

		sp.SetError(err)

		log.Printf("[Sub Wrapper] After serving publication")

		return err
	}
}
