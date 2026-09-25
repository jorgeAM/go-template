package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jorgeAM/go-template/internal/shared/crypto"
)

type contextKey string

const principalContextKey contextKey = "principal"

// Principal is the authenticated caller as carried by the access token. Subject is the JWT "sub"
// claim; modules map it (and any extra Claims) onto their own concepts.
type Principal struct {
	Subject string
	Claims  map[string]any
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorHandler(w, "authentication required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errorHandler(w, "invalid authorization header format")
			return
		}

		jwtToken := parts[1]

		claims, err := crypto.ValidateTokenWithType(jwtToken, "access")
		if err != nil {
			errorHandler(w, fmt.Sprintf("invalid token: %s", err.Error()))
			return
		}

		subject, ok := claims["sub"].(string)
		if !ok || subject == "" {
			errorHandler(w, "invalid token payload")
			return
		}

		principal := &Principal{
			Subject: subject,
			Claims:  claims,
		}

		ctx := context.WithValue(r.Context(), principalContextKey, principal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const refreshTokenMaxAgeSeconds = 30 * 24 * 3600

func SetAuthCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   refreshTokenMaxAgeSeconds,
	})
}

func ClearAuthCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
}

func GetPrincipalFromContext(ctx context.Context) (*Principal, bool) {
	principal, ok := ctx.Value(principalContextKey).(*Principal)
	return principal, ok
}

func errorHandler(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, message)))
}
