package auth

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
)

type service struct {
	UserRepository     users.Repository
	ConsumerRepository consumers.Repository
	S3Repository       s3.Repository
	DBTx               db.SqlxTransaction
}

var _ Service = (*service)(nil)

type NewServiceOpts struct {
	UserRepository     users.Repository
	ConsumerRepository consumers.Repository
	S3Repository       s3.Repository
	DBTx               db.SqlxTransaction
}

func NewService(
	opts NewServiceOpts,
) *service {
	return &service{
		UserRepository:     opts.UserRepository,
		ConsumerRepository: opts.ConsumerRepository,
		S3Repository:       opts.S3Repository,
		DBTx:               opts.DBTx,
	}
}
