package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

var (
	jwtAccessSecret  []byte
	jwtRefreshSecret []byte
	accessTokenTTL   int
	refreshTokenTTL  int
)

func init() {
	jwtAccessSecret = GetEnvBytes("JWT_ACCESS_SECRET", "default-access-secret")
	jwtRefreshSecret = GetEnvBytes("JWT_REFRESH_SECRET", "default-refresh-secret")
	accessTokenTTL = GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15)
	refreshTokenTTL = GetEnvInt("JWT_REFRESH_TTL_DAYS", 7)
}

func GenerateAccessToken(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(accessTokenTTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtAccessSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenTTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtRefreshSecret)
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr, jwtAccessSecret)
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr, jwtRefreshSecret)
}

func validateToken(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
