package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"net/url"
)

type minioClient interface {
	PresignedPostPolicy(ctx context.Context, policy *minio.PostPolicy) (url *url.URL, formData map[string]string, err error)
	GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}
