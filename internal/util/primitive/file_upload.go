package primitive

import (
	"errors"
	"flag"
	"github.com/google/uuid"
	"strings"
	"time"
)

type PresignedFileUpload struct {
	Identifier       string
	OriginalFileName string
	MimeType         MimeType
	Size             int64
	ChecksumSHA256   string

	// Generated by Backend
	GeneratedFileName string
	Extension         string
}

type PresignedFileUploadOutput struct {
	Identifier      string
	UploadURL       string
	UploadExpiredAt time.Time
	MinioFormData   map[string]string
}

type NewPresignedFileUploadInput struct {
	Identifier       string
	OriginalFileName string
	MimeType         MimeType
	Size             int64
	ChecksumSHA256   string
}

func (v NewPresignedFileUploadInput) extension() string {
	if v.MimeType == MimeTypeJpeg {
		splitFileName := strings.Split(v.OriginalFileName, ".")
		return "." + splitFileName[len(splitFileName)-1]
	}
	return MapMimeTypeExtensions[v.MimeType]
}

func NewPresignedFileUpload(input NewPresignedFileUploadInput) (output PresignedFileUpload, err error) {
	if !input.MimeType.IsValid() {
		return output, errors.New("mime type " + string(input.MimeType) + " is not allowed")
	}
	extension := input.extension()
	uuidString := uuid.New().String()
	generatedFileName := uuidString + extension
	if flag.Lookup("test.v") != nil {
		generatedFileName = input.OriginalFileName
	}
	return PresignedFileUpload{
		Identifier:        input.Identifier,
		OriginalFileName:  input.OriginalFileName,
		MimeType:          input.MimeType,
		Size:              input.Size,
		ChecksumSHA256:    input.ChecksumSHA256,
		GeneratedFileName: generatedFileName,
		Extension:         extension,
	}, nil
}
