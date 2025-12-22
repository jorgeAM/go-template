package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "github.com/jorgeAM/go-template/cfg"
	"github.com/jorgeAM/go-template/pkg/log"

	_ "github.com/joho/godotenv/autoload"
)

func startServer(cfg *config.Config, deps *config.Dependencies) error {
	router := buildRouter(cfg, deps)

	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := log.InitDefaultLogger(); err != nil {
		log.Panic(ctx, "error initializing default logger", log.WithError(err))
	}

	log.Info(ctx, "[Config] Loading...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(ctx, "error loading config", log.WithError(err))
	}

	log.Info(ctx, "[Config] Finished")
	log.Info(ctx, "[Dependencies] Building...")

	deps, err := config.BuildDependencies(cfg)
	if err != nil {
		log.Panic(ctx, "error building dependencies", log.WithError(err))
	}

	log.Info(ctx, "[Dependencies] Finished")

	log.Info(ctx, "[App] Initializing")
	go func() {
		log.Info(ctx, fmt.Sprintf("[Server] Listening on %s", cfg.Port))

		if err := startServer(cfg, deps); err != nil {
			log.Panic(ctx, "error starting server", log.WithError(err))
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
}
