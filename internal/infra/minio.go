package infra

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinio(cred conf.ConfigMinio) *minio.Client {
	minioClient, err := minio.New(cred.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cred.AccessID, cred.SecretAccessKey, ""),
		Secure: cred.UseSSL,
	})
	util.Panic(err)

	exist, err := minioClient.BucketExists(context.Background(), cred.PrivateBucket)
	util.Panic(err)

	if !exist {
		err = minioClient.MakeBucket(context.Background(), cred.PrivateBucket, minio.MakeBucketOptions{})
		util.Panic(err)
	}

	log.Info().Msg("initialization minio successfully")
	return minioClient
}
