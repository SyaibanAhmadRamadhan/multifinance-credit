package middleware

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func (m *middleware) StartingOtelTrace(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("starting otel trace").Start(r.Context(), r.URL.Host+r.URL.Path, trace.WithAttributes(
			attribute.String("request.method", r.Method),
			attribute.String("request.user_agent", r.UserAgent()),
			attribute.String("request.content-type", r.Header.Get("Content-Type")),
			attribute.Int64("request.content-length", r.ContentLength),
		))
		defer span.End()

		ctx = context.WithValue(ctx, primitive.SpanIDKey, span.SpanContext().SpanID().String())
		w.Header().Set("X-Request-ID", span.SpanContext().SpanID().String())

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
