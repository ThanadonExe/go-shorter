package services

import (
	"errors"

	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
	"github.com/thanadonexe/go-shorter/internal/utils/base62"
)

type UrlService struct {
	repo ports.UrlRepository
}

func NewUrlService(repo ports.UrlRepository) *UrlService {
	return &UrlService{
		repo: repo,
	}
}

func (s *UrlService) Create(req *domain.UrlCreateRequest) (*domain.UrlResponse, error) {
	id, err := s.repo.NextUrlID()
	if err != nil {
		return nil, err
	}

	req.ID = id
	if req.ShortUrl == "" {
		req.ShortUrl = base62.Encode(id)
	}

	url, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return url.ToUrlResponse(), nil
}

func (s *UrlService) GetById(id uint64) (*domain.UrlResponse, error) {
	url, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return url.ToUrlResponse(), nil
}

func (s *UrlService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *UrlService) GetAll(page, limit int) (*domain.UrlListResponse, error) {
	records := s.repo.Count()
	totalPage := records / limit
	if records%limit != 0 {
		totalPage += 1
	}

	if page > totalPage {
		return nil, errors.New("page exceeded")
	}

	offset := (page - 1) * limit
	urls, err := s.repo.GetAll(offset, limit)
	if err != nil {
		return nil, err
	}

	var data []*domain.UrlResponse
	for _, u := range urls {
		data = append(data, u.ToUrlResponse())
	}

	urlList := &domain.UrlListResponse{
		Data:      data,
		Total:     records,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
	}

	return urlList, nil
}

func (s *UrlService) GetByCode(code string) (*domain.UrlResponse, error) {
	url, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	return url.ToUrlResponse(), nil
}
