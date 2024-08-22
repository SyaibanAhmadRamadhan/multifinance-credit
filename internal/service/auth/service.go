package auth

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3"
)

type service struct {
	userRepository     users.Repository
	consumerRepository consumers.Repository
	limitRepository    limits.Repository
	s3Repository       s3.Repository
	dbTx               db.SqlxTransaction
}

var _ Service = (*service)(nil)

type NewServiceOpts struct {
	UserRepository     users.Repository
	ConsumerRepository consumers.Repository
	LimitRepository    limits.Repository
	S3Repository       s3.Repository
	DBTx               db.SqlxTransaction
}

func NewService(
	opts NewServiceOpts,
) *service {
	return &service{
		userRepository:     opts.UserRepository,
		consumerRepository: opts.ConsumerRepository,
		limitRepository:    opts.LimitRepository,
		s3Repository:       opts.S3Repository,
		dbTx:               opts.DBTx,
	}
}
