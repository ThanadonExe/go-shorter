package ports

import "github.com/thanadonexe/go-shorter/internal/core/domain"

type AuthService interface {
	Login(request *domain.AuthRequest) (*domain.Credential, error)
	Refresh(request *domain.RefreshRequest) (*domain.Credential, error)
}
