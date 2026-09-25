package bootstrap

import (
	"github.com/jorgeAM/go-template/internal/identity"
	"github.com/jorgeAM/go-template/internal/shared/module"
)

func Modules() []module.Module {
	return []module.Module{
		identity.NewModule(),
	}
}
