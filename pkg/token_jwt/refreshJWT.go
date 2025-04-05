package tokenjwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ClaimsRefresh описывает содержимое refresh-токена
type ClaimsRefresh struct {
	UserID    string `json:"user_id"`
	Secret    string `json:"secret"`
	RefreshId string `json:"refresh_id"`
	jwt.RegisteredClaims
}

// GenerateSecureSecret генерирует криптоустойчивую строку
func GenerateSecureSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateRefreshJWT создает refresh-токен
func GenerateRefreshJWT(userID uuid.UUID, secret string, refreshId uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &ClaimsRefresh{
		UserID:    userID.String(),
		Secret:    secret,
		RefreshId: refreshId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return signedToken, nil
}

// DecodeRefreshJWT разбирает и валидирует refresh-токен
func DecodeRefreshJWT(tokenStr string) (*ClaimsRefresh, error) {
	claims := &ClaimsRefresh{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil // ✅ Правильно
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token signature")
	}

	return claims, nil
}
