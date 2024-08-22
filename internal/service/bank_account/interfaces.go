package bank_account

import "context"

type Service interface {
	Creates(ctx context.Context, input CreatesInput) (err error)
	GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error)
}
