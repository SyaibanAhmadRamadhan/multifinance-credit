package presentation

import (
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/conf"
	middlewarecustom "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/middleware"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

type Presenter struct {
	DependencyService *service.Dependency
}

func New(presenter *Presenter) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3002"},
		AllowedHeaders:   []string{"Origin", "Test", "Content-Type", "Accept", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))

	handler(presenter, r)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.GetConfig().AppPort),
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	return server
}

func handler(presenter *Presenter, r *chi.Mux) {
	restApi := restapi.New(presenter.DependencyService)
	middlewareCustom := middlewarecustom.NewMiddleware(presenter.DependencyService.AuthService)
	r.Use(middlewareCustom.StartingOtelTrace)

	r.Group(func(r chi.Router) {
		r.Use(middlewareCustom.AuthUser)
		r.Get("/api/v1/image-private", withOtel(
			restApi.V1ImagePrivateGet,
			WithLogResponseBody(false)),
		)
		r.Post("/api/v1/bank-account", withOtel(
			restApi.V1BankAccountsPost,
		))
		r.Get("/api/v1/bank-account", withOtel(
			restApi.V1BankAccountsGet,
		))
	})

	r.Post("/api/v1/register", withOtel(
		restApi.V1RegisterPost,
		WithLogRequestBody(false)),
	)

	r.Post("/api/v1/login", withOtel(
		restApi.V1LoginPost,
		WithLogResponseBody(false),
		WithLogRequestBody(false)),
	)

	r.Post("/api/v1/refresh-token", withOtel(
		restApi.V1RefreshTokenPost,
		WithLogResponseBody(false),
		WithLogRequestBody(false)),
	)
}
