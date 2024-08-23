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

type GetObjectInput struct {
	ObjectName string
}

type GetObjectOutput struct {
	Object io.ReadCloser
}

type GetPresignedUrlInput struct {
	ObjectName string
	BucketName string
	Expired    time.Duration
}

type GetPresignedUrlOutput struct {
	URL string
}
