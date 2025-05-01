package tracing

import (
	"fmt"

	"je-suis-ici-activitypub/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitJaeger(cfg *config.JaegerConfig) (*tracesdk.TracerProvider, error) {
	if !cfg.Enable {
		return nil, nil
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.URL)))
	if err != nil {
		return nil, fmt.Errorf("fail to create jaeger exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
