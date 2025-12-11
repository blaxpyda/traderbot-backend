package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

// getSecretKey retrieves the secret key for signing JWT tokens.
func getSecretKey() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "default_secret_key"
	}
	return []byte(s)
}

// GenerateToken generates a JWT token with custom claims.
func GenerateJWT(accountID string, ttl time.Duration) (string, error) {
	if accountID == "" {
		return "", errors.New("accountID cannot be empty")
	}
	if ttl <= 0 {
		ttl = time.Hour // Default to 1 hour
	}
	now := time.Now()
	claims := CustomClaims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

// ParseAndValidateJWT parses and validates a JWT token string.
func ParseAndValidateJWT(tokenStr string) (*CustomClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	
	if claims, ok := parsed.Claims.(*CustomClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}