package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := opentracing.GlobalTracer()
		if tracer == nil {
			ctx.Next()
			return
		}
		var span opentracing.Span

		wireCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			span = opentracing.StartSpan(ctx.Request.URL.Path)
		} else {
			span = opentracing.StartSpan(ctx.Request.URL.Path, opentracing.ChildOf(wireCtx))
		}

		defer span.Finish()
		if b := ctx.Request.GetBody; b != nil {
			if body, err := b(); err == nil {
				var buf bytes.Buffer
				if _, err = buf.ReadFrom(body); err == nil {
					span.LogKV("form-body", buf.String())
				}
			}
		}
		if sp, ok := span.Context().(jaeger.SpanContext); ok {
			ctx.Set("root_span", span)
			ctx.Writer.Header().Set("x-request-id", sp.TraceID().String())
			ctx.Writer.Header().Set("x-trace-id", sp.TraceID().String())
			ctx.Writer.Header().Set("X-Span-id", sp.SpanID().String())
		}
		ctx.Next()
	}
}
