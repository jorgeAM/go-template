package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jorgeAM/go-template/internal/bootstrap"
	"github.com/jorgeAM/go-template/internal/platform/log"
	"github.com/jorgeAM/go-template/internal/shared/env"
	"github.com/jorgeAM/go-template/internal/shared/module"

	_ "github.com/joho/godotenv/autoload"
)

func startServer(ctx context.Context, port string, modules []module.Module) error {
	router, err := buildRouter(ctx, modules)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := log.InitDefaultLogger(); err != nil {
		log.Panic(ctx, "error initializing default logger", log.WithError(err))
	}

	port := env.GetEnv("PORT", "8080")

	log.Info(ctx, "[Modules] Initializing")

	modules := bootstrap.Modules()
	for _, m := range modules {
		if err := m.Init(ctx); err != nil {
			log.Panic(ctx, "error initializing module", log.WithString("module", string(m.Name())), log.WithError(err))
		}
	}

	log.Info(ctx, "[Modules] Finished")

	log.Info(ctx, "[App] Initializing")
	go func() {
		log.Info(ctx, "[Server] Listening", log.WithString("port", port))

		if err := startServer(ctx, port, modules); err != nil {
			log.Panic(ctx, "error starting server", log.WithError(err))
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
}
