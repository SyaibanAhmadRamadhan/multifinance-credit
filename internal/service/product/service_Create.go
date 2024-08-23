package product

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"golang.org/x/sync/errgroup"
)

func (s *service) Create(ctx context.Context, input CreateInput) (output CreateOutput, err error) {
	var erg errgroup.Group

	erg.Go(func() (err error) {
		outputPresignedUrl, err := s.s3Repository.CreatePresignedUrl(ctx, s3.CreatePresignedUrlInput{
			BucketName: conf.GetConfig().Minio.PrivateBucket,
			Path:       input.Image.GeneratedFileName,
			MimeType:   string(input.Image.MimeType),
			Checksum:   input.Image.ChecksumSHA256,
		})
		if err != nil {
			return tracer.Error(err)
		}

		output = CreateOutput{
			Image: primitive.PresignedFileUploadOutput{
				Identifier:      input.Image.Identifier,
				UploadURL:       outputPresignedUrl.URL,
				UploadExpiredAt: outputPresignedUrl.ExpiredAt,
				MinioFormData:   outputPresignedUrl.MinioFormData,
			},
		}
		return nil
	})

	err = s.dbTx.DoTransaction(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}, func(tx *db.SqlxWrapper) error {

		err = s.productRepository.Creates(ctx, products.CreatesInput{
			Transaction: tx,
			MerchantID:  input.MerchantID,
			Items: []products.CreatesInputItem{
				{
					Name:  input.Name,
					Image: input.Image.GeneratedFileName,
					Qty:   input.Qty,
					Price: input.Price,
				},
			},
		})
		if err != nil {
			return tracer.Error(err)
		}

		if err = erg.Wait(); err != nil {
			return tracer.Error(err)
		}

		return nil
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}
