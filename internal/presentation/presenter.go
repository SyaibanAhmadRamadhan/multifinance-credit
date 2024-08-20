package presentation

import (
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

const SpanIDKey = "span_id"

type Presenter struct{}

func New(presenter *Presenter) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	handler(presenter, r)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.GetConfig().AppPort),
		Handler: r,
	}
	return server
}

func handler(presenter *Presenter, r *chi.Mux) {
	restApi := restapi.New()
	r.Get("/hello", WithOtel(restApi.HelloWorld))
}
