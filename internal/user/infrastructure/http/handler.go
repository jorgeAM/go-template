package http

import (
	"encoding/json"
	"net/http"

	config "github.com/jorgeAM/go-template/cfg"
	"github.com/jorgeAM/go-template/internal/user/application"
	"github.com/jorgeAM/go-template/pkg/http/response"
)

func CreateUser(_ *config.Config, deps *config.Dependencies) http.HandlerFunc {
	srv := application.NewCreateUser(deps.UserRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		var body application.CreateUserCommand
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.BadRequest(w, "BAD_REQUEST", err.Error())
			return
		}

		if err := srv.Exec(r.Context(), &body); err != nil {
			response.InternalServerErr(w, "INTERNAL_ERROR", err.Error())
			return
		}

		response.OK(w, "ok")
	}
}
