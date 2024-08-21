package consumers

import "context"

type Repository interface {
	Create(ctx context.Context, input CreateInput) (output CreateOutput, err error)
	CheckExisting(ctx context.Context, input CheckExistingInput) (output CheckExistingOutput, err error)
	Get(ctx context.Context, input GetInput) (output GetOutput, err error)
}
