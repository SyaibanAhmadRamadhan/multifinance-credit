package minio

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/url"
)

func (r *repository) GetPresignedUrl(ctx context.Context, input s3.GetPresignedUrlInput) (output s3.GetPresignedUrlOutput, err error) {
	ctx, span := r.tracer.Start(ctx, "minio s3: Get Presigned Url", trace.WithAttributes(
		attribute.String("bucket_name", input.BucketName),
		attribute.String("object_name", input.ObjectName),
	))
	defer span.End()

	params := make(url.Values)

	getPresignedOutput, err := r.client.PresignedGetObject(ctx, input.BucketName, input.ObjectName, input.Expired, params)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	output = s3.GetPresignedUrlOutput{
		URL: getPresignedOutput.String(),
	}
	return
}
