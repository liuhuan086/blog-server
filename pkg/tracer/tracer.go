package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second,
			LogSpans:            true,
			LocalAgentHostPort:  agentHostPort,
		},
	}

	trace, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(trace)
	return trace, closer, nil
}
