package product_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/product"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

func Test_service_GetAll(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()

	conf.Init()
	mockProductRepository := products.NewMockRepository(mock)
	mockS3Repository := s3.NewMockRepository(mock)

	s := product.NewService(product.NewServiceOpts{
		ProductRepository: mockProductRepository,
		S3Repository:      mockS3Repository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(2)
		expectedImageGetPresignedUrl := faker.URL()
		expectedPagination := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}
		expectedInput := product.GetAllInput{
			MerchantID: null.IntFrom(rand.Int63()),
			IDs:        []int64{1, 2},
			Pagination: expectedPagination,
		}
		expectedProductRepoImage1 := faker.UUIDDigit()
		expectedProductRepoImage2 := faker.UUIDDigit()
		expectedOutput := product.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPagination, expectedTotalData),
			Items: []product.GetAllOutputItem{
				{
					ID:         rand.Int63(),
					MerchantID: expectedInput.MerchantID.Int64,
					Name:       faker.Name(),
					Image:      expectedImageGetPresignedUrl,
					Qty:        rand.Int63(),
					Price:      rand.Float64(),
				},
				{
					ID:         rand.Int63(),
					MerchantID: expectedInput.MerchantID.Int64,
					Name:       faker.Name(),
					Image:      expectedImageGetPresignedUrl,
					Qty:        rand.Int63(),
					Price:      rand.Float64(),
				},
			},
		}

		mockS3Repository.EXPECT().
			GetPresignedUrl(ctx, s3.GetPresignedUrlInput{
				ObjectName: expectedProductRepoImage1,
				BucketName: conf.GetConfig().Minio.PrivateBucket,
				Expired:    5 * time.Minute,
			}).
			Return(s3.GetPresignedUrlOutput{
				URL: expectedImageGetPresignedUrl,
			}, nil)
		mockS3Repository.EXPECT().
			GetPresignedUrl(ctx, s3.GetPresignedUrlInput{
				ObjectName: expectedProductRepoImage2,
				BucketName: conf.GetConfig().Minio.PrivateBucket,
				Expired:    5 * time.Minute,
			}).
			Return(s3.GetPresignedUrlOutput{
				URL: expectedImageGetPresignedUrl,
			}, nil)

		mockProductRepository.EXPECT().
			GetAll(ctx, products.GetAllInput{
				MerchantID: expectedInput.MerchantID,
				IDs:        expectedInput.IDs,
				Pagination: expectedPagination,
			}).
			Return(products.GetAllOutput{
				Pagination: pagination.CreatePaginationOutput(expectedPagination, expectedTotalData),
				Items: []products.GetAllOutputItem{
					{
						ID:         expectedOutput.Items[0].ID,
						MerchantID: expectedOutput.Items[0].MerchantID,
						Image:      expectedProductRepoImage1,
						Name:       expectedOutput.Items[0].Name,
						Qty:        expectedOutput.Items[0].Qty,
						Price:      expectedOutput.Items[0].Price,
					},
					{
						ID:         expectedOutput.Items[1].ID,
						MerchantID: expectedOutput.Items[1].MerchantID,
						Image:      expectedProductRepoImage2,
						Name:       expectedOutput.Items[1].Name,
						Qty:        expectedOutput.Items[1].Qty,
						Price:      expectedOutput.Items[1].Price,
					},
				},
			}, nil)

		output, err := s.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
