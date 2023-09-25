package provider

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// collector url
func newTraceProvider(url string) (*sdktrace.TracerProvider, error) {
	expr, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
		),
	)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(expr),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(os.Getenv("APP_NAME")),
			attribute.String("environment", os.Getenv("QUORA_ENV")),
		)),
	), nil
}

func ProvideOtel() (*sdktrace.TracerProvider, trace.Tracer, error) {
	traceProvider, err := newTraceProvider(os.Getenv("JAEGER_EXPORTER_URL"))
	if err != nil {
		return nil, nil, err
	}

	otel.SetTracerProvider(traceProvider)

	tp := traceProvider.Tracer("quora-clone")

	return traceProvider, tp, nil
}
