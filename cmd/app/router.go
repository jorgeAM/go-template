package main

import (
	"net/http"
	"time"

	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	config "github.com/jorgeAM/go-template/cfg"
	userHandler "github.com/jorgeAM/go-template/internal/user/infrastructure/http"
	"github.com/jorgeAM/go-template/pkg/http/handler"
	"github.com/jorgeAM/go-template/pkg/http/middleware"
)

func buildRouter(cfg *config.Config, deps *config.Dependencies) http.Handler {
	router := chi.NewRouter()

	httpin_integration.UseGochiURLParam("path", chi.URLParam)

	router.Use(
		middleware.RequestID,
		middleware.Logger(middleware.WithIgnoreRoutes("/health")),
		middleware.Recover,
		middleware.RealIP,
		middleware.CORS(middleware.DefaultCORSOptions),
		middleware.ResponseHeader("Content-Type", "application/json"),
		middleware.ResponseHeader("Accept", "application/json"),
		middleware.Timeout(15*time.Second),
	)

	router.Get("/health", handler.HealthCheck)

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser(cfg, deps))
			r.Get("/{id}", userHandler.GetUser(cfg, deps))
		})
	})

	return router
}
