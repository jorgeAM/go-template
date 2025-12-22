package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jorgeAM/go-template/pkg/crypto"
	"github.com/jorgeAM/go-template/pkg/env"
)

type contextKey string

const USER_CONTEXT_KEY contextKey = "user"

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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

		userID, idOk := claims["sub"].(string)
		email, emailOk := claims["email"].(string)
		if !idOk || !emailOk {
			errorHandler(w, "invalid token payload")
			return
		}

		userInfo := &UserInfo{
			ID:    userID,
			Email: email,
		}

		ctx := context.WithValue(r.Context(), USER_CONTEXT_KEY, userInfo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   env.GetEnv("APP_ENV", "local") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days to match JWT creation
	})
}

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   env.GetEnv("APP_ENV", "local") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
}

func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(USER_CONTEXT_KEY).(*UserInfo)
	return user, ok
}

func errorHandler(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, message)))
}
