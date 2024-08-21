package s3

import (
	"io"
	"time"
)

type CreatePresignedUrlInput struct {
	BucketName string
	Path       string
	MimeType   string
	Checksum   string
}

type CreatePresignedUrlOutput struct {
	URL           string
	ExpiredAt     time.Time
	MinioFormData map[string]string
}

type GetPrivateObjectInput struct {
	ObjectName string
}

type GetPrivateObjectOutput struct {
	Object io.ReadCloser
}
