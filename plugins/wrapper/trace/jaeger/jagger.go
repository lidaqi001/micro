package jaeger

import (
	"context"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"github.com/lidaqi001/micro/examples/config"
	"time"
)

// 创建一个jaeger Tracer
func NewTracer(servicename string) (opentracing.Tracer, io.Closer, error) {
	addr := config.TRACE_ADDR
	//traceIp := getTraceIp(addr)

	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			//CollectorEndpoint:   "http://" + traceIp + ":14268/api/traces",
			//LocalAgentHostPort:  "192.168.1.145:6831",
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}

// 创建空的上下文, 生成追踪 span
func GetTraceClientCtxAndSpan() (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "call services")
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	defer span.Finish()

	// 注入 opentracing textmap 到空的上下文用于追踪
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)

	return span, ctx
}

// 从微服务上下文中获取追踪信息
func GetTraceServiceSpan(ctx context.Context) opentracing.Span {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	// 创建新的 Span 并将其绑定到微服务上下文
	return opentracing.StartSpan("call Server", opentracing.ChildOf(wireContext))
}
