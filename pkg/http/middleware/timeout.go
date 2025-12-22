package middleware

import (
	"context"
	"net/http"
	"time"
)

var timeoutHeader = http.CanonicalHeaderKey("X-Timeout")

func Timeout(defaultTimeout time.Duration) Middleware {
	if defaultTimeout == 0 {
		defaultTimeout = 100 * time.Millisecond
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeout := defaultTimeout

			if timeoutStr := r.Header.Get(timeoutHeader); timeoutStr != "" {
				d, err := time.ParseDuration(timeoutStr)
				if err == nil {
					timeout = d
				}
			}

			if timeout < 10*time.Millisecond {
				timeout = 10 * time.Millisecond
			}

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()

			r.Header.Set(timeoutHeader, timeout.String())

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
