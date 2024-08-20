package util

import "context"

type CloseFn func(ctx context.Context) (err error)
