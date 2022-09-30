package conf

import (
	"context"
	"main/util"
	"os"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

func RegisterTracerProvider() {
	backend := os.Getenv("BACKEND")

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("api"),
	)

	if backend == "XRAY" {
		withXRayExporter(res)
	} else if backend == "JAEGER" {
		withJaegerExporter(res)
	} else {
		util.PanicLog(nil, "OpenTelemetry BACKEND not set")
	}
}

func withXRayExporter(res *resource.Resource) {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:4317"
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		util.PanicLog(err, "failed to create new OTLP trace exporter for X-Ray")
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)

	// If you extend trace beyond a single process, use context propagation: https://opentelemetry.io/docs/instrumentation/go/manual/#propagators-and-context
	//otel.SetTextMapPropagator(xray.Propagator{})
}

func withJaegerExporter(res *resource.Resource) {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:4317"
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		util.PanicLog(err, "failed to create new OTLP trace exporter for Jaeger")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(tp)
}
