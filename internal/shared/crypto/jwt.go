package crypto

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Sentinel errors are safe to return to clients; the underlying jwt library error never is.
var (
	ErrTokenMalformed = errors.New("token is malformed or has an invalid signature")
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenType      = errors.New("unexpected token type")
	ErrTokenIssuer    = errors.New("unexpected token issuer")
)

func GenerateJWT(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ValidateToken(jwtToken string) (jwt.Claims, error) {
	return parse(jwtToken)
}

func ValidateTokenWithType(jwtToken string, expectedType string) (jwt.MapClaims, error) {
	claims, err := parse(jwtToken)
	if err != nil {
		return nil, err
	}

	tokenType, err := ExtractTokenType(claims)
	if err != nil {
		return nil, err
	}

	if tokenType != expectedType {
		return nil, ErrTokenType
	}

	issuer, ok := claims["iss"].(string)
	if !ok || issuer != os.Getenv("JWT_ISSUER") {
		return nil, ErrTokenIssuer
	}

	return claims, nil
}

func ExtractTokenType(claims jwt.MapClaims) (string, error) {
	tokenType, ok := claims["type"].(string)
	if !ok {
		return "", ErrTokenType
	}

	return tokenType, nil
}

func parse(jwtToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrTokenMalformed
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrTokenMalformed
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenMalformed
	}

	return claims, nil
}
