package products

import "context"

type Repository interface {
	Get(ctx context.Context, input GetInput) (output GetOutput, err error)
	GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error)
	Creates(ctx context.Context, input CreatesInput) (err error)
	Updates(ctx context.Context, input UpdatesInput) (err error)
}
