package main

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/infra"
)

func main() {
	conf.Init()

	infra.NewMinio(conf.GetConfig().Minio)
	infra.NewOtel(conf.GetConfig().OpenTelemetry)
	infra.NewMysql(conf.GetConfig().DatabaseDSN)
}
