package middleware

import (
	"net/http"
	"strings"
)

func ResponseHeader(key, value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.TrimSpace(w.Header().Get(key)) == "" {
				w.Header().Set(key, value)
			}

			next.ServeHTTP(w, r)
		})
	}
}
