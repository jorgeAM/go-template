package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Recover(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}
