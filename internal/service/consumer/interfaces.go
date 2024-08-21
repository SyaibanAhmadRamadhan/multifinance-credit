package consumer

import "context"

type Service interface {
	GetPrivateImage(ctx context.Context, input GetPrivateImageInput) (output GetPrivateImageOutput, err error)
}
