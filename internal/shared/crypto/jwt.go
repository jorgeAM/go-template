package crypto

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ValidateToken(jwtToken string) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims, nil
}

func ValidateTokenWithType(jwtToken string, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	tokenType, err := ExtractTokenType(claims)
	if err != nil {
		return nil, err
	}

	if tokenType != expectedType {
		return nil, errors.New("wrong token type")
	}

	issuer, ok := claims["iss"].(string)
	if !ok || issuer != os.Getenv("JWT_ISSUER") {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}

func ExtractTokenType(claims jwt.MapClaims) (string, error) {
	tokenType, ok := claims["type"].(string)
	if !ok {
		return "", errors.New("missing or invalid token type")
	}

	return tokenType, nil
}
