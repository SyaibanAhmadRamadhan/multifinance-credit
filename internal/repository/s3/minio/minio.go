package minio

import (
	"github.com/jonboulle/clockwork"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	client minioClient
	tracer trace.Tracer
	attrs  []attribute.KeyValue
	clock  clockwork.Clock
}

func NewRepository(client minioClient, clock clockwork.Clock) *repository {
	tp := otel.GetTracerProvider()
	return &repository{
		client: client,
		tracer: tp.Tracer("minio-tracer", trace.WithInstrumentationVersion("v1.0.0")),
		attrs: []attribute.KeyValue{
			attribute.String("minio-library-version", "v7"),
		},
		clock: clock,
	}
}
