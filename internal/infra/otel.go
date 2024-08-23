package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"time"
)

func NewOtel(cred conf.ConfigOpenTelemetry) util.CloseFn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cred.Username, cred.Password)))
	traceCli := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithHeaders(map[string]string{
			"Authorization": authHeader,
		}),
		otlptracegrpc.WithEndpoint(cred.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	traceExp, err := otlptrace.New(ctx, traceCli)
	util.Panic(err)

	traceProvide, closeFnTracer, err := newOtelProvider(cred.TracerName).start(traceExp)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed initializing the tracer provider")
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvide)

	log.Info().Msg("initialization otlp successfully")
	return closeFnTracer
}

type otelProvider struct {
	name     string
	exporter trace.SpanExporter
}

func newOtelProvider(name string) *otelProvider {
	return &otelProvider{
		name:     name,
		exporter: nil,
	}
}

func (t *otelProvider) start(exp trace.SpanExporter) (*trace.TracerProvider, util.CloseFn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(t.name),
		),
	)

	if err != nil {
		err = fmt.Errorf("failed to created resource: %w", err)
		return nil, nil, err
	}

	t.exporter = exp
	bsp := trace.NewBatchSpanProcessor(t.exporter)

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	closeFn := func(ctx context.Context) (err error) {
		log.Info().Msg("starting shutdown export and provider")
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = t.exporter.Shutdown(ctx); err != nil {
			return err
		}

		if err = provider.Shutdown(ctx); err != nil {
			return err
		}
		
		log.Info().Msg("shutdown export and provider successfully")

		return
	}
	return provider, closeFn, nil
}
