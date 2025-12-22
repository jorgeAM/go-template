package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}
