package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func newExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New()
}

func newTracerProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	ctx := context.Background()

	exp, err := newExporter()
	if err != nil {
		log.Fatalf("failed to initialize exporter %v", err)
	}

	tp := newTracerProvider(exp)

	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer("ExampleService")

	ctx, span := tracer.Start(ctx, "main func")
	defer span.End()

	sayHello(ctx)

	fmt.Println("End main")
}

func sayHello(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "sayHello")
	defer span.End()

	fmt.Println("Hello")
}
