package interceptor

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const traceIDKey = "x-trace-id"

// ServerTracingInterceptor intercepts and traces gRPC requests
func ServerTracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIDKey, spanContext.TraceID().String()))

		header := metadata.New(map[string]string{traceIDKey: spanContext.TraceID().String()})
		err := grpc.SendHeader(ctx, header)
		if err != nil {
			return nil, err
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
	}

	return res, err
}
