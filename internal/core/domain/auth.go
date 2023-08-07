package domain

import "github.com/golang-jwt/jwt/v5"

type JWTAccessTokenClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTRefreshTokenClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Credential struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
