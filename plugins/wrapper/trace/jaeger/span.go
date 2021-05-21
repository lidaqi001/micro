package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

type span struct {
	span opentracing.Span
}

type Span interface {
	SetTopic(topic interface{})
	SetError(error interface{})
	SetHeader(header interface{})
	SetPayload(payload interface{})
	SetRequest(request interface{})
	SetResponse(response interface{}, err error, opts ...bool)
}

func NewSpan(ctx context.Context) Span {
	// 从微服务上下文中关联链路信息
	// 创建新的 Span 并将其绑定到微服务上下文
	sp := GetTraceServiceSpan(ctx)
	return &span{span: sp}
}

func (s span) SetTopic(req interface{}) {
	s.span.SetTag("topic", req)
}

func (s span) SetHeader(hea interface{}) {
	s.span.SetTag("header", hea)
}

func (s span) SetPayload(hea interface{}) {
	s.span.SetTag("payload", hea)
}

func (s span) SetRequest(req interface{}) {
	s.span.SetTag("request", req)
}

func (s span) SetError(err interface{}) {
	if err != nil {
		s.span.SetTag("err", err)
	}
	s.span.Finish()
}

func (s span) SetResponse(rsp interface{}, err error, opts ...bool) {

	s.span.SetTag("response", rsp)

	s.SetError(err)

	if len(opts) > 0 {
		if opts[0] == false {
			return
		}
	}
	// 服务有可能相互调用，不确定在哪个服务是最后一次响应，所以不设置span的finish标签
	//s.span.Finish()
}
