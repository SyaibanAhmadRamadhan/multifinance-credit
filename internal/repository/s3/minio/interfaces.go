package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"net/url"
	"time"
)

type minioClient interface {
	PresignedPostPolicy(ctx context.Context, policy *minio.PostPolicy) (url *url.URL, formData map[string]string, err error)
	GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error)
}
