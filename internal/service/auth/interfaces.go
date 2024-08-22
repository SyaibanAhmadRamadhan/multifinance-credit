package auth

import "context"

type Service interface {
	Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error)
	Login(ctx context.Context, input LoginInput) (output LoginOutput, err error)
	VerifyToken(ctx context.Context, input VerifyTokenInput) (output VerifyTokenOutput, err error)
}
