package users

import "context"

type Repository interface {
	CheckExisting(ctx context.Context, input CheckExistingInput) (output CheckExistingOutput, err error)
	Create(ctx context.Context, input CreateInput) (output CreateOutput, err error)
	Get(ctx context.Context, input GetInput) (output GetOutput, err error)
}
