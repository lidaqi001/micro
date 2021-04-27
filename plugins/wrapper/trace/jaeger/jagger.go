package jaeger

import (
	"context"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"os"
	"time"
)

const DEFAULT_TRACE_IP = "192.168.1.146"

// NewTracer 创建一个jaeger Tracer
func NewTracer(servicename string, addr string, ip string) (opentracing.Tracer, io.Closer, error) {
	traceIp := getTraceIp(ip)

	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   "http://" + traceIp + ":14268/api/traces",
			//LocalAgentHostPort:  "192.168.1.145:6831",
		},
	}

	sender, err := jaeger.NewUDPTransport(traceIp+":"+addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, err
}

func getTraceIp(ip string) string {
	// 传入IP就使用该参数
	traceIp := ip
	if empty(traceIp) {
		// 传入IP为空，则从env获取
		traceIp = os.Getenv("MICRO_TRACE_IP")
	}

	if empty(traceIp) {
		// 否则使用默认IP
		traceIp = DEFAULT_TRACE_IP
	}

	if empty(traceIp) {
		log.Println("traceIp:::" + traceIp)
		log.Println("ip:::MICRO_TRACE_IP:::" + os.Getenv("MICRO_TRACE_IP"))
		log.Println("ip:::DEFAULT_TRACE_IP:::" + DEFAULT_TRACE_IP)
		panic("Trace server ip is null!")
	}

	log.Println("traceIp:::" + traceIp)
	return traceIp
}

func empty(ip string) bool {
	return len(ip) == 0
}

func GetTraceClientCtxAndSpan() (opentracing.Span, context.Context) {

	// 创建空的上下文, 生成追踪 span
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "call Service")
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

func GetTraceServiceSpan(ctx *context.Context, req interface{}, rsp interface{}) opentracing.Span {
	// 从微服务上下文中获取追踪信息
	md, ok := metadata.FromContext(*ctx)
	if !ok {
		md = make(map[string]string)
	}
	var sp opentracing.Span
	wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	// 创建新的 Span 并将其绑定到微服务上下文
	//sp = opentracing.StartSpan("SayHello", opentracing.ChildOf(wireContext))
	sp = opentracing.StartSpan("call Server", opentracing.ChildOf(wireContext))

	// 记录请求
	sp.SetTag("req", req)

	// 同时记录响应
	if rsp != nil {
		SpanSetResponse(sp, rsp)
	}

	return sp
}

func SpanSetResponse(sp opentracing.Span, rsp interface{}) {
	// 记录响应
	sp.SetTag("resp", rsp)
	// 在函数返回 stop span 之前，统计函数执行时间
	sp.Finish()
}
