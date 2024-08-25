package product_test

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

func Test_service_Creates(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()

	conf.Init()
	mockProductRepository := products.NewMockRepository(mock)
	mockDBTx := db.NewMockSqlxTransaction(mock)
	mockS3Repository := s3.NewMockRepository(mock)

	s := product.NewService(product.NewServiceOpts{
		ProductRepository: mockProductRepository,
		DBTx:              mockDBTx,
		S3Repository:      mockS3Repository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := product.CreateInput{
			MerchantID: rand.Int63(),
			Name:       faker.Name(),
			Image: primitive.PresignedFileUpload{
				Identifier:        faker.UUIDDigit(),
				OriginalFileName:  "image.png",
				MimeType:          "image/png",
				Size:              5000,
				ChecksumSHA256:    faker.UUIDDigit(),
				GeneratedFileName: faker.UUIDDigit(),
				Extension:         "png",
			},
			Qty:   rand.Int63(),
			Price: rand.Float64(),
		}
		expectedPresignedUrlOutput := primitive.PresignedFileUploadOutput{
			Identifier:      expectedInput.Image.Identifier,
			UploadURL:       faker.URL(),
			UploadExpiredAt: time.Now().UTC(),
			MinioFormData:   make(map[string]string),
		}
		mockS3Repository.EXPECT().
			CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
				BucketName: conf.GetConfig().Minio.PrivateBucket,
				Path:       expectedInput.Image.GeneratedFileName,
				MimeType:   string(expectedInput.Image.MimeType),
				Checksum:   expectedInput.Image.ChecksumSHA256,
			}).
			Return(s3.CreatePresignedUrlOutput{
				URL:           expectedPresignedUrlOutput.UploadURL,
				ExpiredAt:     expectedPresignedUrlOutput.UploadExpiredAt,
				MinioFormData: expectedPresignedUrlOutput.MinioFormData,
			}, nil)

		mockDBTx.EXPECT().
			DoTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
				gomock.Any()).
			DoAndReturn(func(ctx context.Context, tx *sql.TxOptions, fn func(tx db.Rdbms) error) error {
				mockSqlxWrapper := db.NewRdbms(nil)
				mockProductRepository.EXPECT().
					Creates(ctx, products.CreatesInput{
						Transaction: mockSqlxWrapper,
						MerchantID:  expectedInput.MerchantID,
						Items: []products.CreatesInputItem{
							{
								Name:  expectedInput.Name,
								Image: expectedInput.Image.GeneratedFileName,
								Qty:   expectedInput.Qty,
								Price: expectedInput.Price,
							},
						},
					}).Return(nil)
				return fn(mockSqlxWrapper)
			}).Return(nil)

		output, err := s.Create(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedPresignedUrlOutput, output.Image)
	})
}
