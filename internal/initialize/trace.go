package initialize

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitTrace() opentracing.Tracer {
	if opentracing.IsGlobalTracerRegistered() {
		return opentracing.GlobalTracer()
	} else {
		cfg := config.Configuration{
			ServiceName: "api.douyacun.com",
			Sampler: &config.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: "127.0.0.1:6831",
			},
			Headers: &jaeger.HeadersConfig{
				JaegerDebugHeader:        "x-debug-id",
				JaegerBaggageHeader:      "x-baggage",
				TraceContextHeaderName:   "x-trace-id",
				TraceBaggageHeaderPrefix: "x-ctx",
			},
		}
		tracer, _, err := cfg.NewTracer()
		if err != nil {
			panic(errors.Wrapf(err, "Error: cannot init tracer"))
		}
		opentracing.SetGlobalTracer(tracer)
		return tracer
	}
}
