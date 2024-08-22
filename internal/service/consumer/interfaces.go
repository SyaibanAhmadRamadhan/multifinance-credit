package consumer

import "context"

type Service interface {
	GetPrivateImage(ctx context.Context, input GetPrivateImageInput) (output GetPrivateImageOutput, err error)
	Get(ctx context.Context, input GetInput) (output GetOutput, err error)
}
