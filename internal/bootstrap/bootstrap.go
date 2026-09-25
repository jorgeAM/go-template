// Package bootstrap holds the one module list, so cmd/app (server) and
// cmd/migrate (migrations) wire a new module in from a single place.
package bootstrap

import (
	"github.com/jorgeAM/go-template/internal/shared/module"
	"github.com/jorgeAM/go-template/internal/user"
)

func Modules() []module.Module {
	return []module.Module{
		user.NewModule(),
	}
}
