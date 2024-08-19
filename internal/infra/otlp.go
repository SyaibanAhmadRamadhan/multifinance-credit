package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"time"
)

func NewOTLP(cred conf.ConfigOpenTelemetry) *otlptrace.Exporter {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cred.Username, cred.Password)))
	traceCli := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithHeaders(map[string]string{
			"Authorization": authHeader,
		}),
		otlptracegrpc.WithEndpoint(cred.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	traceExp, err := otlptrace.New(ctx, traceCli)
	util.Panic(err)

	log.Info().Msg("initialization otlp successfully")
	return traceExp
}
