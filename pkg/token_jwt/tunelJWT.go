package tokenjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TunnelClaims struct {
	jwt.RegisteredClaims
}

// GenerateTunnelToken создает JWT токен для ПК с TTL = 5 лет
func GenerateTunnelToken(computerID string) (string, error) {
	now := time.Now()
	exp := now.Add(5 * 365 * 24 * time.Hour) // 5 лет

	claims := TunnelClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        computerID,               // jti: ID туннеля/ПК
			Subject:   computerID,               // sub: ID ПК
			Issuer:    "zentas",                 // iss
			Audience:  []string{"zic-protocol"}, // aud
			IssuedAt:  jwt.NewNumericDate(now),  // iat
			ExpiresAt: jwt.NewNumericDate(exp),  // exp
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

// ParseTunnelToken разбирает и валидирует JWT токен
func ParseTunnelToken(tokenStr string) (*TunnelClaims, error) {
	claims := &TunnelClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
