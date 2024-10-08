package main

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/infra"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/jonboulle/clockwork"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	conf.Init()

	minioClient := infra.NewMinio(conf.GetConfig().Minio)
	otel := infra.NewOtel(conf.GetConfig().OpenTelemetry)
	sqlxDB, dbClose := infra.NewMysql(conf.GetConfig().DatabaseDSN)
	clockWork := clockwork.NewRealClock()

	server := presentation.New(&presentation.Presenter{
		DependencyService: service.NewDependency(service.NewDependencyOpts{
			MinioClient: minioClient,
			SqlxDB:      sqlxDB,
			Clock:       clockWork,
		}),
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-quit
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	if err := otel(ctx); err != nil {
		panic(err)
	}

	if err := dbClose(ctx); err != nil {
		panic(err)
	}

	log.Info().Msg("Server gracefully stopped")
}
