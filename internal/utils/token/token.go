package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
)

func CreateAccessToken(user *domain.User, secret string, expire int64) (string, error) {
	expireAt := time.Now().Add(time.Duration(expire) * time.Hour)
	claims := &domain.JWTAccessTokenClaims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tkn, nil
}

func CreateRefreshToken(user *domain.User, secret string, expire int64) (string, error) {
	expireAt := time.Now().Add(time.Duration(expire) * time.Hour)
	claims := &domain.JWTRefreshTokenClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tkn, nil
}

func ParseToken(token, secret string) (*domain.JWTAccessTokenClaims, error) {
	tkn, err := jwt.ParseWithClaims(token, &domain.JWTAccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(*domain.JWTAccessTokenClaims)
	if !ok || !tkn.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidateToken(token, secret string) (bool, error) {
	_, err := ParseToken(token, secret)

	if err != nil {
		return false, err
	}

	return true, nil
}
