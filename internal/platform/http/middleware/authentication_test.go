package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/jorgeAM/go-template/internal/shared/crypto"
)

func TestAuthenticate(t *testing.T) {
	const issuer = "go-template-test"

	t.Setenv("JWT_KEY", "test-secret")
	t.Setenv("JWT_ISSUER", issuer)

	mint := func(t *testing.T, claims jwt.MapClaims) string {
		t.Helper()
		claims["iss"] = issuer
		if _, ok := claims["exp"]; !ok {
			claims["exp"] = time.Now().Add(time.Hour).Unix()
		}
		token, err := crypto.GenerateJWT(claims)
		assert.NoError(t, err)
		return token
	}

	accessClaims := func() jwt.MapClaims {
		return jwt.MapClaims{"sub": "account-1", "type": "access", "role": "member"}
	}

	tests := []struct {
		name          string
		authHeader    func(t *testing.T) string
		wantStatus    int
		wantPrincipal bool
	}{
		{
			name:       "missing authorization header is rejected",
			authHeader: func(*testing.T) string { return "" },
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "non-bearer authorization header is rejected",
			authHeader: func(t *testing.T) string { return "Basic " + mint(t, accessClaims()) },
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "expired token is rejected",
			authHeader: func(t *testing.T) string {
				c := accessClaims()
				c["exp"] = time.Now().Add(-time.Hour).Unix()
				return "Bearer " + mint(t, c)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "refresh token is rejected",
			authHeader: func(t *testing.T) string {
				c := accessClaims()
				c["type"] = "refresh"
				return "Bearer " + mint(t, c)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "token without a subject is rejected",
			authHeader: func(t *testing.T) string {
				c := accessClaims()
				delete(c, "sub")
				return "Bearer " + mint(t, c)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:          "valid token passes and populates the principal",
			authHeader:    func(t *testing.T) string { return "Bearer " + mint(t, accessClaims()) },
			wantStatus:    http.StatusOK,
			wantPrincipal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Principal
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got, _ = GetPrincipalFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/whatever", nil)
			if h := tt.authHeader(t); h != "" {
				req.Header.Set("Authorization", h)
			}

			rec := httptest.NewRecorder()
			Authenticate(next).ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if !tt.wantPrincipal {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, "account-1", got.Subject)
			assert.Equal(t, "member", got.Claims["role"])
		})
	}
}

func TestErrorHandler_EmitsValidJSON(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	errorHandler(rec, `invalid token: contains a " quote and \ backslash`)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body struct {
		Message string `json:"message"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, `invalid token: contains a " quote and \ backslash`, body.Message)
}
