package middleware

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/auth"
	"go.opentelemetry.io/otel"
)

var tracerName = "middleware"
var otelTracer = otel.Tracer(tracerName)

type middleware struct {
	authService auth.Service
}

func NewMiddleware(authService auth.Service) *middleware {
	return &middleware{authService: authService}
}
