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
	"github.com/spf13/cobra"
	"net/http"
	"os/signal"
	"syscall"
)

var restApiCmd = &cobra.Command{
	Use:   "rest-api",
	Short: "run rest api",
	Run: func(cmd *cobra.Command, args []string) {
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

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := server.Shutdown(context.TODO()); err != nil {
			panic(err)
		}

		if err := otel(context.TODO()); err != nil {
			panic(err)
		}

		if err := dbClose(context.TODO()); err != nil {
			panic(err)
		}

		log.Info().Msg("Shutdown complete. Exiting.")
		return
	},
}
