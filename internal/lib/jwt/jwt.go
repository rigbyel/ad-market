package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rigbyel/ad-market/internal/models"
)

type UserClaims struct {
	jwt.RegisteredClaims
	ID    int64
	Login string
}

// NewToken creates new JWT token for given user
func NewToken(user models.User, secret string, duration time.Duration) (string, error) {
	const op = "lib.jwt.NewToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration))},
		ID:               user.Id,
		Login:            user.Login,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

// GetTokenClaims parses claims from token string
func GetTokenClaims(tokenString, secret string) (UserClaims, error) {
	const op = "lib.jwt.GetTokenClaims"
	var claims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: wrong token signing method", op)
		}

		return []byte(secret), nil
	})

	if err != nil {
		return UserClaims{}, fmt.Errorf("%s: %w", op, err)
	}
	if !token.Valid {
		return UserClaims{}, fmt.Errorf("invalid token")
	}

	return claims, nil
}
