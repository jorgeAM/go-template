package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jorgeAM/go-template/internal/platform/http/response"
	"github.com/jorgeAM/go-template/internal/user/application/command"
	"github.com/jorgeAM/go-template/internal/user/application/query"
	"github.com/jorgeAM/go-template/internal/user/domain"
)

func Register(r chi.Router, userRepository domain.UserRepository) {
	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/", CreateUser(userRepository))
		r.Get("/{id}", GetUser(userRepository))
	})
}

func CreateUser(userRepository domain.UserRepository) http.HandlerFunc {
	srv := command.NewCreateUser(userRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		var body command.CreateUserCommand
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

func GetUser(userRepository domain.UserRepository) http.HandlerFunc {
	srv := query.NewGetUser(userRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			response.BadRequest(w, "BAD_REQUEST", "user id is required")
			return
		}

		res, err := srv.Exec(r.Context(), userID)
		if err != nil {
			response.InternalServerErr(w, "INTERNAL_ERROR", err.Error())
			return
		}

		response.OK(w, res)
	}
}
