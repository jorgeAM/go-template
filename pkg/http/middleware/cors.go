package middleware

import (
	"github.com/go-chi/cors"
)

type CORSOptions = cors.Options

var DefaultCORSOptions = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

func CORS(options CORSOptions) Middleware {
	return cors.Handler(options)
}
