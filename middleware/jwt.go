package middleware

import (
	"context"
	"net/http"
	"strings"

	"thugcorp.io/final_bot/utils"
)

type contextKey string

const (
	ContextKeyAccountID contextKey = "account_id"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimSpace(authHeader[len("Bearer "):])
		claims, err := utils.ParseAndValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyAccountID, claims.AccountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetJWTClaims extracts custom claims from a request context.
func GetJWTClaims(r *http.Request) *utils.CustomClaims {
	accountID, ok := r.Context().Value(ContextKeyAccountID).(string)
	if !ok {
		return nil
	}
	return &utils.CustomClaims{
		AccountID: accountID,
	}
}