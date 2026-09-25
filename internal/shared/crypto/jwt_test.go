package crypto

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidateTokenWithType(t *testing.T) {
	const issuer = "go-template-test"

	t.Setenv("JWT_KEY", "test-secret")
	t.Setenv("JWT_ISSUER", issuer)

	base := func() jwt.MapClaims {
		return jwt.MapClaims{
			"sub":  "account-1",
			"type": "access",
			"iss":  issuer,
			"exp":  time.Now().Add(time.Hour).Unix(),
		}
	}

	mint := func(t *testing.T, claims jwt.MapClaims) string {
		t.Helper()
		raw, err := GenerateJWT(claims)
		assert.NoError(t, err)
		return raw
	}

	tests := []struct {
		name    string
		raw     func(t *testing.T) string
		wantErr error
	}{
		{
			name: "valid access token",
			raw:  func(t *testing.T) string { return mint(t, base()) },
		},
		{
			name: "expired token",
			raw: func(t *testing.T) string {
				c := base()
				c["exp"] = time.Now().Add(-time.Hour).Unix()
				return mint(t, c)
			},
			wantErr: ErrTokenExpired,
		},
		{
			name: "wrong token type",
			raw: func(t *testing.T) string {
				c := base()
				c["type"] = "refresh"
				return mint(t, c)
			},
			wantErr: ErrTokenType,
		},
		{
			name: "missing token type",
			raw: func(t *testing.T) string {
				c := base()
				delete(c, "type")
				return mint(t, c)
			},
			wantErr: ErrTokenType,
		},
		{
			name: "unknown issuer",
			raw: func(t *testing.T) string {
				c := base()
				c["iss"] = "someone-else"
				return mint(t, c)
			},
			wantErr: ErrTokenIssuer,
		},
		{
			name:    "garbage input",
			raw:     func(*testing.T) string { return "not-a-jwt" },
			wantErr: ErrTokenMalformed,
		},
		{
			name: "signed with the wrong key",
			raw: func(t *testing.T) string {
				t.Setenv("JWT_KEY", "attacker-secret")
				raw := mint(t, base())
				t.Setenv("JWT_KEY", "test-secret")
				return raw
			},
			wantErr: ErrTokenMalformed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateTokenWithType(tt.raw(t), "access")

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, claims)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "account-1", claims["sub"])
		})
	}
}
