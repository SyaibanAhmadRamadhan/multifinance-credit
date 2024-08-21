package minio_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3/minio"
	"github.com/go-faker/faker/v4"
	"github.com/jonboulle/clockwork"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_repository_GetPrivateObject(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()
	ctx := context.Background()
	mockClientMinio := minio.NewMockminioClient(mock)
	mockClockWork := clockwork.NewFakeClock()

	m := minio.NewRepository(mockClientMinio, mockClockWork)

	t.Run("should return correct object", func(t *testing.T) {
		expectedInput := s3.GetPrivateObjectInput{
			ObjectName: faker.Name(),
		}

		mockClientMinio.EXPECT().
			GetObject(gomock.Any(), "private", expectedInput.ObjectName, miniogo.GetObjectOptions{}).
			Return(nil, nil)

		_, err := m.GetPrivateObject(ctx, expectedInput)
		require.NoError(t, err)
	})
}
