package minio_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3/minio"
	"github.com/go-faker/faker/v4"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/url"
	"testing"
	"time"
)

func Test_repository_GetPresignedUrl(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mockClientMinio := minio.NewMockminioClient(mock)
	mockClockWork := clockwork.NewFakeClock()

	m := minio.NewRepository(mockClientMinio, mockClockWork)

	t.Run("should be return correct", func(t *testing.T) {
		params := make(url.Values)
		expectedInput := s3.GetPresignedUrlInput{
			ObjectName: faker.Name(),
			BucketName: faker.Name(),
			Expired:    time.Hour * 24,
		}
		expectedUrl := &url.URL{
			Scheme: "https",
			Host:   "minio.aliyuncs.com",
		}

		mockClientMinio.EXPECT().
			PresignedGetObject(gomock.Any(), expectedInput.BucketName, expectedInput.ObjectName, expectedInput.Expired, params).
			Return(expectedUrl, nil)

		output, err := m.GetPresignedUrl(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedUrl.String(), output.URL)
	})
}
