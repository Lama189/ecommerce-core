package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
)

type contextKey string

const ctxUserIDKey contextKey = "userID"

func AuthMiddleware(jwtManager user.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, http.StatusUnauthorized, "authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			payload, err := jwtManager.VerifyToken(parts[1], "access")
			if err != nil {
				sendError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserIDKey, payload.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
