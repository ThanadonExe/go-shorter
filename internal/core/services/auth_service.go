package services

import (
	"fmt"

	"github.com/thanadonexe/go-shorter/internal/config"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
	"github.com/thanadonexe/go-shorter/internal/utils/hash"
	"github.com/thanadonexe/go-shorter/internal/utils/token"
)

type AuthService struct {
	repo ports.UserRepository
}

func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Login(request *domain.AuthRequest) (*domain.Credential, error) {
	user, err := s.repo.GetByEmail(request.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if ok := hash.CheckBCryptHash(request.Password, user.Password); !ok {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, err := token.CreateAccessToken(user, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiredIn)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token")
	}

	refreshToken, err := token.CreateRefreshToken(user, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiredIn)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token")
	}

	res := &domain.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (s *AuthService) Refresh(request *domain.RefreshRequest) (*domain.Credential, error) {
	claims, err := token.ParseToken(request.RefreshToken, config.AppConfig.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	id := claims.ID
	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	accessToken, err := token.CreateAccessToken(user, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiredIn)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token")
	}

	refreshToken, err := token.CreateRefreshToken(user, config.AppConfig.JWTSecret, 24)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token")
	}

	res := &domain.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
