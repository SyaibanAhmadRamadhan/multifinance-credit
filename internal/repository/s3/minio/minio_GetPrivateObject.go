package minio

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (r *repository) GetPrivateObject(ctx context.Context, input s3.GetPrivateObjectInput) (output s3.GetPrivateObjectOutput, err error) {
	ctx, span := r.tracer.Start(ctx, "minio s3: Get Private Object", trace.WithAttributes(
		attribute.String("bucket_name", conf.GetConfig().Minio.PrivateBucket),
		attribute.String("object_name", input.ObjectName),
	))
	defer span.End()

	object, err := r.client.GetObject(ctx, conf.GetConfig().Minio.PrivateBucket, input.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	_, err = object.Stat()
	if err != nil {
		err = errors.Join(err, s3.ErrIsBadRequest)
		span.RecordError(err)
		return output, tracer.Error(err)
	}

	output = s3.GetPrivateObjectOutput{
		Object: object,
	}
	return
}
