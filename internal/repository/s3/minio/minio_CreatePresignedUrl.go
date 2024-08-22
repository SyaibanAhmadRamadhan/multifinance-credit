package minio

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func (r *repository) CreatePresignedUrl(ctx context.Context, input s3.CreatePresignedUrlInput) (output s3.CreatePresignedUrlOutput, err error) {
	time.Sleep(5 * time.Second)
	expiredAt := r.clock.Now().UTC().Add(5 * time.Minute)

	ctx, span := r.tracer.Start(ctx, "minio s3 - Create Presigned Url", trace.WithAttributes(
		attribute.String("bucket_name", input.BucketName),
		attribute.String("path_name", input.Path),
		attribute.String("expired", expiredAt.Format(time.DateTime)),
		attribute.String("mime_type", input.MimeType),
	))
	defer span.End()

	policy := minio.NewPostPolicy()
	policy.SetChecksum(minio.NewChecksum(minio.ChecksumSHA256, []byte(input.Checksum)))
	err = policy.SetKey(input.Path)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	err = policy.SetExpires(expiredAt)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	err = policy.SetContentLengthRange(1024, 2048*2048)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	err = policy.SetBucket(input.BucketName)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	outputPresignedPostPolicy, formData, err := r.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		return output, tracer.Error(err)
	}

	output = s3.CreatePresignedUrlOutput{
		URL:           outputPresignedPostPolicy.String(),
		ExpiredAt:     expiredAt,
		MinioFormData: formData,
	}
	return
}
