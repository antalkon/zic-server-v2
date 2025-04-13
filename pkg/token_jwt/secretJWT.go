package tokenjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ClaimsPc описывает содержимое PC-токена
type ClaimsPc struct {
	PcID   string `json:"pc_id"`
	Secret string `json:"secret"`
	jwt.RegisteredClaims
}

// GeneratePcJWT создает JWT для PC с заданным pcID и secret
func GeneratePcJWT(pcID uuid.UUID, secret string) (string, error) {
	expiration := time.Now().Add(5 * 365 * 24 * time.Hour) // 30 дней

	claims := &ClaimsPc{
		PcID:   pcID.String(),
		Secret: secret,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

// DecodePcJWT парсит JWT токен и возвращает ClaimsPc
func DecodePcJWT(tokenStr string) (*ClaimsPc, error) {
	claims := &ClaimsPc{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse PC token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid PC token signature")
	}

	return claims, nil
}
