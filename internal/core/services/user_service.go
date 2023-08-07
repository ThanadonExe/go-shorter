package services

import (
	"errors"

	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
	"github.com/thanadonexe/go-shorter/internal/utils/hash"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(request *domain.UserCreateRequest) (*domain.UserResponse, error) {
	pass, err := hash.HashBCrypt(request.Password)
	if err != nil {
		return nil, err
	}

	request.Password = pass

	_, err = s.repo.GetByEmail(request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	user, err := s.repo.Create(request)
	if err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (s *UserService) GetById(id int) (*domain.UserResponse, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (s *UserService) Update(id int, request *domain.UserUpdateRequest) (*domain.UserResponse, error) {
	pass, err := hash.HashBCrypt(request.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	request.Password = pass

	user.ApplyUpdate(request)

	if err = s.repo.Update(user); err != nil {
		return nil, err
	}

	return user.ToUserResponse(), nil
}

func (s *UserService) Delete(id int) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *UserService) GetAll(page, limit int) (*domain.UserListResponse, error) {
	records := s.repo.Count()
	totalPage := records / limit
	if records%limit != 0 {
		totalPage += 1
	}

	if page > totalPage {
		return nil, errors.New("page exceeded")
	}

	offset := (page - 1) * limit
	users, err := s.repo.GetAll(offset, limit)
	if err != nil {
		return nil, err
	}

	var data []*domain.UserResponse
	for _, u := range users {
		data = append(data, u.ToUserResponse())
	}

	userList := &domain.UserListResponse{
		Data:      data,
		Total:     records,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
	}

	return userList, nil
}

func (s *UserService) GetByEmail(email string) (*domain.UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user.ToUserResponse(), err
}
