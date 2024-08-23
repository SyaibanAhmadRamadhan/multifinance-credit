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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_repository_GetPrivateObject(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()
	conf.Init()
	ctx := context.Background()
	mockClockWork := clockwork.NewFakeClock()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Last-Modified", "Wed, 21 Oct 2015 07:28:00 GMT")
		w.Header().Set("Content-Length", "5")

		w.Write([]byte("12345"))
	}))
	defer srv.Close()

	clnt, err := miniogo.New(srv.Listener.Addr().String(), &miniogo.Options{
		Region: "us-east-1",
	})
	require.NoError(t, err)

	m := minio.NewRepository(clnt, mockClockWork)

	t.Run("should return correct object", func(t *testing.T) {
		expectedInput := s3.GetObjectInput{
			ObjectName: faker.Name(),
		}

		output, err := m.GetObject(ctx, expectedInput)
		require.NoError(t, err)

		buf, err := io.ReadAll(output.Object)
		require.NoError(t, err)
		require.Equal(t, "12345", string(buf))
	})
}
