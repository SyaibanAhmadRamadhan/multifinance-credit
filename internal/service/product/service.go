package product

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
)

type service struct {
	productRepository products.Repository
	s3Repository      s3.Repository
	dbTx              db.SqlxTransaction
}

var _ Service = (*service)(nil)

type NewServiceOpts struct {
	ProductRepository products.Repository
	S3Repository      s3.Repository
	DBTx              db.SqlxTransaction
}

func NewService(
	opts NewServiceOpts,
) *service {
	return &service{
		productRepository: opts.ProductRepository,
		s3Repository:      opts.S3Repository,
		dbTx:              opts.DBTx,
	}
}
