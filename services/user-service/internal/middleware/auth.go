package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SovetkanB/payflow/user-service/internal/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response{
					Success: false,
					Data:    "authorization is required",
				})
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := auth.ValidateToken(tokenStr, secret)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response{
					Success: false,
					Data:    "authorization is required",
				})
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) *auth.Claims {
	v := ctx.Value(ClaimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*auth.Claims)
	return claims
}
