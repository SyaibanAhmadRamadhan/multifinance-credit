package consumer

import (
	"github.com/guregu/null/v5"
	"io"
)

type GetPrivateImageInput struct {
	UserID      int64
	ConsumerID  null.Int
	ImageKtp    null.Bool
	ImageSelfie null.Bool
}

type GetPrivateImageOutput struct {
	Object io.ReadCloser
}

type GetInput struct {
	UserID     null.Int
	ConsumerID null.Int
}

type GetOutput struct {
	ID        int64
	UserID    int64
	FullName  string
	LegalName string
}
