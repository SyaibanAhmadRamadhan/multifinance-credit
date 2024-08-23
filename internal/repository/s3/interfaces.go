package s3

import "context"

type Repository interface {
	CreatePresignedUrl(ctx context.Context, input CreatePresignedUrlInput) (output CreatePresignedUrlOutput, err error)
	GetObject(ctx context.Context, input GetObjectInput) (output GetObjectOutput, err error)
	GetPresignedUrl(ctx context.Context, input GetPresignedUrlInput) (output GetPresignedUrlOutput, err error)
}
