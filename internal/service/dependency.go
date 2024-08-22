package service

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	minio_repository "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/s3/minio"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/consumer"
	"github.com/jmoiron/sqlx"
	"github.com/jonboulle/clockwork"
	"github.com/minio/minio-go/v7"
)

type Dependency struct {
	AuthService     auth.Service
	ConsumerService consumer.Service
}

type NewDependencyOpts struct {
	MinioClient *minio.Client
	SqlxDB      *sqlx.DB
	Clock       clockwork.Clock
}

func NewDependency(opts NewDependencyOpts) *Dependency {
	sqlxWrapper := db.NewSqlxWrapper(opts.SqlxDB)
	sqlxTransaction := db.NewSqlxTransaction(opts.SqlxDB)

	// REPOSITORY LAYER
	userRepository := users.NewRepository(sqlxWrapper)
	consumerRepository := consumers.NewRepository(sqlxWrapper)
	minioRepository := minio_repository.NewRepository(opts.MinioClient, opts.Clock)

	// SERVICE LAYER
	authService := auth.NewService(auth.NewServiceOpts{
		UserRepository:     userRepository,
		ConsumerRepository: consumerRepository,
		S3Repository:       minioRepository,
		DBTx:               sqlxTransaction,
	})
	consumerService := consumer.NewService(consumer.NewServiceOpts{
		UserRepository:     userRepository,
		ConsumerRepository: consumerRepository,
		S3Repository:       minioRepository,
		DBTx:               sqlxTransaction,
	})
	return &Dependency{
		AuthService:     authService,
		ConsumerService: consumerService,
	}
}
