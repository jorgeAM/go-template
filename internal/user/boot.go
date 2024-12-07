package user

import (
	"github.com/go-chi/chi/v5"
	config "github.com/jorgeAM/base-api/cfg"
)

type Bootable struct {
	cfg  *config.Config
	deps *config.Dependencies
}

func Boot(cfg *config.Config, deps *config.Dependencies) (*Bootable, error) {
	return &Bootable{
		cfg:  cfg,
		deps: deps,
	}, nil
}

func (b *Bootable) BuildRoutes(router chi.Router) {
	router.Post("/", createUser(b.cfg, b.deps))
}
