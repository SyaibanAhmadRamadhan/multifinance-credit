package s3

import "context"

type Repository interface {
	CreatePresignedUrl(ctx context.Context, input CreatePresignedUrlInput) (output CreatePresignedUrlOutput, err error)
	GetPrivateObject(ctx context.Context, input GetPrivateObjectInput) (output GetPrivateObjectOutput, err error)
}
