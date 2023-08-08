package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/thanadonexe/go-shorter/internal/adapters/caches"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
)

type UrlRepositoryMock struct {
	urls  []*domain.Url
	cache ports.CacheRepository
}

func NewUrlRepositoryMock() *UrlRepositoryMock {
	cache, _ := caches.NewRedisCacheMock()
	urls := []*domain.Url{
		{ID: 1, ShortUrl: "L1", FullUrl: "http://www.google.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, ShortUrl: "L2", FullUrl: "http://www.google.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, ShortUrl: "L3", FullUrl: "http://www.google.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	return &UrlRepositoryMock{
		urls:  urls,
		cache: cache,
	}
}

func (r *UrlRepositoryMock) NextUrlID() (uint64, error) {
	var id uint64

	for _, u := range r.urls {
		if u.ID > id {
			id = u.ID
		}
	}

	id += 1
	return id, nil
}

func (r *UrlRepositoryMock) Create(req *domain.UrlCreateRequest) (*domain.Url, error) {
	url := domain.Url{
		ID:        req.ID,
		ShortUrl:  req.ShortUrl,
		FullUrl:   req.FullUrl,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &url, nil
}

func (r *UrlRepositoryMock) GetById(id uint64) (*domain.Url, error) {
	var url domain.Url
	for _, u := range r.urls {
		if u.ID == id && u.DeletedAt.IsZero() {
			url = *u
		}
	}

	if url == (domain.Url{}) {
		return nil, sql.ErrNoRows
	}

	return &url, nil
}

func (r *UrlRepositoryMock) Delete(id uint64) error {
	for _, u := range r.urls {
		if u.ID == id {
			u.DeletedAt = time.Now()
		}
	}

	return nil
}

func (r *UrlRepositoryMock) GetAll(offset, limit int) ([]*domain.Url, error) {
	if len(r.urls) <= 0 {
		return nil, sql.ErrNoRows
	}

	return r.urls, nil
}

func (r *UrlRepositoryMock) GetByCode(code string) (*domain.Url, error) {
	var url domain.Url
	cacheKey := fmt.Sprintf("%s_%s", UrlCachePrefix, code)
	err := r.cache.Get(cacheKey, &url)

	for _, u := range r.urls {
		if u.ShortUrl == code && u.DeletedAt.IsZero() {
			url = *u
		}
	}

	if url == (domain.Url{}) {
		return nil, sql.ErrNoRows
	}

	_ = r.cache.Set(cacheKey, &url, UrlCacheDuration)

	return &url, err
}

func (r *UrlRepositoryMock) Count() int {
	return len(r.urls)
}
