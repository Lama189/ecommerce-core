package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTManager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  []byte(secretKey),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *JWTManager) GenerateTokenPair(payload user.TokenPayload) (string, string, error) {
	accessToken, err := m.generateToken(payload, TokenTypeAccess, m.accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := m.generateToken(payload, TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) VerifyToken(tokenStr string, expectedTYpe string) (*user.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing methid: %v", t.Header["alg"])
		}

		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != expectedTYpe {
		return nil, ErrInvalidType
	}

	return &user.TokenPayload{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func (m *JWTManager) generateToken(payload user.TokenPayload, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := userClaims{
		UserID: payload.UserID,
		Role:   payload.Role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}
