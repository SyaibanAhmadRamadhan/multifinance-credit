package consumer_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func Test_service_GetPrivateImage(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockS3Repository := s3.NewMockRepository(mock)
	mockConsumerRepository := consumers.NewMockRepository(mock)
	ctx := context.Background()

	s := consumer.NewService(consumer.NewServiceOpts{
		ConsumerRepository: mockConsumerRepository,
		S3Repository:       mockS3Repository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := consumer.GetPrivateImageInput{
			UserID:     rand.Int63(),
			ConsumerID: null.IntFrom(rand.Int63()),
			ImageKtp:   null.BoolFrom(true),
		}
		expectedIoRead := io.NopCloser(strings.NewReader("dummy"))
		expectedOutput := consumer.GetPrivateImageOutput{
			Object: expectedIoRead,
		}

		mockConsumerRepository.EXPECT().
			Get(ctx, consumers.GetInput{
				ID:     expectedInput.ConsumerID,
				UserID: null.IntFrom(expectedInput.UserID),
			}).
			Return(consumers.GetOutput{
				PhotoKTP:    "image-ktp.png",
				PhotoSelfie: "image-selfie.jpg",
			}, nil)

		mockS3Repository.EXPECT().
			GetPrivateObject(ctx, s3.GetPrivateObjectInput{
				ObjectName: "image-ktp.png",
			}).
			Return(s3.GetPrivateObjectOutput{
				Object: expectedIoRead,
			}, nil)

		output, err := s.GetPrivateImage(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return correct with empty image selfie", func(t *testing.T) {
		expectedInput := consumer.GetPrivateImageInput{
			UserID:      rand.Int63(),
			ConsumerID:  null.IntFrom(rand.Int63()),
			ImageSelfie: null.BoolFrom(true),
		}
		expectedOutput := consumer.GetPrivateImageOutput{
			Object: nil,
		}

		mockConsumerRepository.EXPECT().
			Get(ctx, consumers.GetInput{
				ID:     expectedInput.ConsumerID,
				UserID: null.IntFrom(expectedInput.UserID),
			}).
			Return(consumers.GetOutput{
				PhotoKTP:    "image-ktp.png",
				PhotoSelfie: "",
			}, nil)

		output, err := s.GetPrivateImage(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error consumer not found", func(t *testing.T) {
		expectedInput := consumer.GetPrivateImageInput{
			UserID:     rand.Int63(),
			ConsumerID: null.IntFrom(rand.Int63()),
			ImageKtp:   null.BoolFrom(true),
		}
		mockConsumerRepository.EXPECT().
			Get(ctx, consumers.GetInput{
				ID:     expectedInput.ConsumerID,
				UserID: null.IntFrom(expectedInput.UserID),
			}).
			Return(consumers.GetOutput{}, datastore.ErrRecordNotFound)

		output, err := s.GetPrivateImage(ctx, expectedInput)
		require.Error(t, err)
		require.ErrorIs(t, err, consumer.ErrConsumerNotFound)
		require.Empty(t, output)
	})
}
