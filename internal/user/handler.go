package user

import (
	"encoding/json"
	"net/http"

	config "github.com/jorgeAM/base-api/cfg"
	"github.com/jorgeAM/base-api/internal/platform/http/response"
	"github.com/jorgeAM/base-api/internal/user/application"
)

func createUser(_ *config.Config, deps *config.Dependencies) http.HandlerFunc {
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
