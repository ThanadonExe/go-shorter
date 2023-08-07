package ports

import "github.com/thanadonexe/go-shorter/internal/core/domain"

type UserService interface {
	Create(user *domain.UserCreateRequest) (*domain.UserResponse, error)
	GetById(id int) (*domain.UserResponse, error)
	Update(id int, request *domain.UserUpdateRequest) (*domain.UserResponse, error)
	Delete(id int) error
	GetAll(page, limit int) (*domain.UserListResponse, error)
	GetByEmail(email string) (*domain.UserResponse, error)
}

type UserRepository interface {
	Create(user *domain.UserCreateRequest) (*domain.User, error)
	GetById(id int) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	GetAll(offset, limit int) ([]*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Count() int
}
