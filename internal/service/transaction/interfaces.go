package transaction

import "context"

type Service interface {
	Create(ctx context.Context, input CreateInput) (output CreateOutput, err error)
}
