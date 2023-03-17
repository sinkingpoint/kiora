package tracing

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const name = "kiora"

func Tracer() trace.Tracer {
	return otel.Tracer(name)
}

type TracingConfiguration struct {
	ServiceName    string
	ExporterType   string
	DestinationURL string
}

func DefaultTracingConfiguration() TracingConfiguration {
	return TracingConfiguration{
		ServiceName:  "kiora",
		ExporterType: "console",
	}
}

func newTracerProvider(config TracingConfiguration, exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	if exp == nil {
		return nil, nil
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNameKey.String(config.ServiceName),
		),
	)

	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	), nil
}

func newSpanExporter(config TracingConfiguration) (sdktrace.SpanExporter, error) {
	switch config.ExporterType {
	case "noop":
		return nil, nil
	case "console":
		return stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
		)
	case "jaeger":
		if config.DestinationURL != "" {
			return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.DestinationURL)))
		}

		return jaeger.New(jaeger.WithCollectorEndpoint())
	default:
		return nil, fmt.Errorf("invalid exporter: %q", config.ExporterType)
	}
}

func InitTracing(config TracingConfiguration) (*sdktrace.TracerProvider, error) {
	exporter, err := newSpanExporter(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create span exporter")
	}

	provider, err := newTracerProvider(config, exporter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tracer provider")
	}

	if provider != nil {
		otel.SetTracerProvider(provider)
	}

	return provider, nil
}
