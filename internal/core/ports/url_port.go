package ports

import "github.com/thanadonexe/go-shorter/internal/core/domain"

type UrlService interface {
	Create(url *domain.UrlCreateRequest) (*domain.UrlResponse, error)
	GetById(id uint64) (*domain.UrlResponse, error)
	Delete(id uint64) error
	GetAll(page, limit int) (*domain.UrlListResponse, error)
	GetByCode(code string) (*domain.UrlResponse, error)
}

type UrlRepository interface {
	Create(url *domain.UrlCreateRequest) (*domain.Url, error)
	GetById(id uint64) (*domain.Url, error)
	Delete(id uint64) error
	GetAll(offset, limit int) ([]*domain.Url, error)
	GetByCode(code string) (*domain.Url, error)
	NextUrlID() (uint64, error)
	Count() int
}
