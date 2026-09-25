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
