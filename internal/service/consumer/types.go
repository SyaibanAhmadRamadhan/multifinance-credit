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
