package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jorgeAM/go-template/pkg/log"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func wrapResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	// Default to 200 OK unless changed
	return &statusResponseWriter{w, http.StatusOK}
}

func Logger(opts ...Option) Middleware {
	options := &loggerOptions{
		ignoredPaths: make(map[string]struct{}),
	}

	options = options.apply(opts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := options.ignoredPaths[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			requestID := middleware.GetReqID(r.Context())
			wrappedWriter := wrapResponseWriter(w)

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			log.Info(
				r.Context(),
				"incoming request",
				log.WithString("method", r.Method),
				log.WithString("path", r.URL.Path),
				log.WithString("remote_ip", r.RemoteAddr),
				log.WithString("user_agent", r.UserAgent()),
				log.WithString("request_id", requestID),
				log.WithInt("status_code", wrappedWriter.statusCode),
				log.WithDuration("duration", duration),
			)
		})
	}
}
