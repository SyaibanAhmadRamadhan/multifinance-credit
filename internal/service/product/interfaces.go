package product

import "context"

type Service interface {
	Create(ctx context.Context, input CreateInput) (output CreateOutput, err error)
	GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error)
}
