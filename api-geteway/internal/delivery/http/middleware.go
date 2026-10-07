package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Lama189/ecommerce-core/api-geteway/internal/infrastructure/jwt"
	"github.com/google/uuid"
)

type contextKey string

const (
	ctxUserIDKey   contextKey = "userID"
	ctxUserRoleKey contextKey = "userRole"
)

func AuthMiddleware(validator jwt.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				SendError(w, http.StatusUnauthorized, "authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				SendError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			claims, err := validator.ValidateAccessToken(parts[1])
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					SendError(w, http.StatusUnauthorized, "token has expired")
				default:
					SendError(w, http.StatusUnauthorized, "invalid or malformed token")
				}
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ctxUserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxUserIDKey).(uuid.UUID)
	return id, ok
}

func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ctxUserRoleKey).(string)
	return role, ok
}
