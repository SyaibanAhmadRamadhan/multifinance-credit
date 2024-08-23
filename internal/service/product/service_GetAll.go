package product

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"time"
)

func (s *service) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	bankAccountsOutput, err := s.productRepository.GetAll(ctx, products.GetAllInput{
		MerchantID: input.MerchantID,
		IDs:        input.IDs,
		Pagination: input.Pagination,
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	output = GetAllOutput{
		Pagination: bankAccountsOutput.Pagination,
		Items:      make([]GetAllOutputItem, 0),
	}

	for _, item := range bankAccountsOutput.Items {
		imageUrl, err := s.s3Repository.GetPresignedUrl(ctx, s3.GetPresignedUrlInput{
			ObjectName: item.Image,
			BucketName: conf.GetConfig().Minio.PrivateBucket,
			Expired:    5 * time.Minute,
		})
		if err != nil {
			return output, tracer.Error(err)
		}

		output.Items = append(output.Items, GetAllOutputItem{
			ID:         item.ID,
			MerchantID: item.MerchantID,
			Name:       item.Name,
			Image:      imageUrl.URL,
			Qty:        item.Qty,
			Price:      item.Price,
		})
	}
	return
}
