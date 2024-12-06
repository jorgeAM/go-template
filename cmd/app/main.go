package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jorgeAM/base-api/internal/platform/log"
)

func main() {
	logger, err := log.NewZapLogger("base-api", "local")
	if err != nil {
		panic(err)
	}

	logger.Info("[App] Initializing")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
}
