package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	config "github.com/jorgeAM/go-template/cfg"
	"github.com/jorgeAM/go-template/internal/platform/http/handler"
	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/user"
)

func startServer(cfg *config.Config, deps *config.Dependencies) error {
	router := chi.NewRouter()

	router.Get("/health", handler.HealthCheck)

	usersBoot, err := user.Boot(cfg, deps)
	if err != nil {
		return err
	}

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", usersBoot.BuildRoutes)
	})

	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)

}

func main() {
	logger, err := log.NewZapLogger("go-template", "local")
	if err != nil {
		panic(err)
	}

	logger.Info("[Config] Loading...")

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger.Info("[Config] Finished")
	logger.Info("[Dependencies] Building...")

	deps, err := config.BuildDependencies(cfg)
	if err != nil {
		panic(err)
	}

	logger.Info("[Dependencies] Finished")

	logger.Info("[App] Initializing")
	go func() {
		logger.Info(fmt.Sprintf("[Server] Listening on %s", cfg.Port))

		if err := startServer(cfg, deps); err != nil {
			logger.Panic("error starting server")
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
}
