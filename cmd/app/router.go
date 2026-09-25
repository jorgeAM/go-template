package main

import (
	"context"
	"net/http"
	"time"

	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/jorgeAM/go-template/internal/platform/http/handler"
	"github.com/jorgeAM/go-template/internal/platform/http/middleware"
	"github.com/jorgeAM/go-template/internal/shared/module"
)

func buildRouter(ctx context.Context, modules []module.Module) (http.Handler, error) {
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

	for _, m := range modules {
		if err := m.RegisterHttp(ctx, router); err != nil {
			return nil, err
		}
	}

	return router, nil
}
