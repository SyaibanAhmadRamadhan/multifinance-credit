package minio_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3/minio"
	"github.com/go-faker/faker/v4"
	"github.com/jonboulle/clockwork"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/url"
	"testing"
	"time"
)

func Test_repository_CreatePresignedUrl(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()
	mockClientMinio := minio.NewMockminioClient(mock)
	mockClockWork := clockwork.NewFakeClock()

	m := minio.NewRepository(mockClientMinio, mockClockWork)

	t.Run("should be return correct", func(t *testing.T) {
		expectedUrl, err := url.Parse(faker.URL())
		require.NoError(t, err)
		expectedInput := s3.CreatePresignedUrlInput{
			BucketName: faker.Name(),
			Path:       "/bucket/image.png",
			MimeType:   "image/jpeg",
			Checksum:   faker.UUIDDigit(),
		}
		expiredAt := mockClockWork.Now().UTC().Add(5 * time.Minute)
		expectedPolicy := miniogo.NewPostPolicy()
		expectedPolicy.SetChecksum(miniogo.NewChecksum(miniogo.ChecksumSHA256, []byte(expectedInput.Checksum)))
		require.NoError(t, expectedPolicy.SetKey(expectedInput.Path))
		require.NoError(t, expectedPolicy.SetExpires(expiredAt))
		require.NoError(t, expectedPolicy.SetContentLengthRange(1024, 2048*2048))
		require.NoError(t, expectedPolicy.SetBucket(expectedInput.BucketName))

		mockClientMinio.EXPECT().
			PresignedPostPolicy(gomock.Any(), expectedPolicy).
			Return(expectedUrl, nil, nil)

		expectedOutput := s3.CreatePresignedUrlOutput{
			URL:       expectedUrl.String(),
			ExpiredAt: expiredAt,
		}

		output, err := m.CreatePresignedUrl(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
