package consumer

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/guregu/null/v5"
)

func (s *service) GetPrivateImage(ctx context.Context, input GetPrivateImageInput) (output GetPrivateImageOutput, err error) {
	userOutput, err := s.ConsumerRepository.Get(ctx, consumers.GetInput{
		ID:     input.ConsumerID,
		UserID: null.IntFrom(input.UserID),
	})
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			err = ErrConsumerNotFound
		}
		return output, tracer.Error(err)
	}

	objectName := ""
	if input.ImageKtp.Bool {
		objectName = userOutput.PhotoKTP
	} else {
		objectName = userOutput.PhotoSelfie
	}

	if objectName == "" {
		return
	}

	privateObjectOutput, err := s.S3Repository.GetPrivateObject(ctx, s3.GetPrivateObjectInput{
		ObjectName: objectName,
	})
	if err != nil {
		return output, err
	}

	output = GetPrivateImageOutput{
		Object: privateObjectOutput.Object,
	}
	return
}
