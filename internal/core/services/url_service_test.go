package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanadonexe/go-shorter/internal/adapters/repositories"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/services"
	"github.com/thanadonexe/go-shorter/internal/utils/base62"
)

func TestGetURLByIDNotFound(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	_, err := service.GetById(9999)

	assert.Error(t, err, "should error")
}

func TestGetURLByIDFound(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	url, err := service.GetById(1)

	assert.NoError(t, err, "should not error")
	assert.NotEqual(t, nil, url, "should not nil")
}

func TestGetAll(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	urls, err := service.GetAll(1, 10)

	assert.NoError(t, err, "should not error")
	assert.Equal(t, 3, len(urls.Data), "data len should be 3")
}

func TestCreateURLWithNoCustomURL(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	id := uint64(4)
	shortUrl := ""
	fullUrl := "http://www.example.com"
	expectedShortUrl := base62.Encode(id)

	req := domain.UrlCreateRequest{
		ID:       id,
		ShortUrl: shortUrl,
		FullUrl:  fullUrl,
	}

	url, err := service.Create(&req)
	assert.NoError(t, err, "should not error")
	assert.Equal(t, id, url.ID, fmt.Sprintf("id must be %d", id))
	assert.Equal(t, "/"+expectedShortUrl, url.ShortUrl, fmt.Sprintf("shortUrl must be %s", expectedShortUrl))
	assert.Equal(t, fullUrl, url.FullUrl, fmt.Sprintf("fullUrl must be %s", fullUrl))
}

func TestCreateURLWithCustomURL(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	id := uint64(4)
	shortUrl := "ex-com"
	fullUrl := "http://www.example.com"

	req := domain.UrlCreateRequest{
		ID:       id,
		ShortUrl: shortUrl,
		FullUrl:  fullUrl,
	}

	url, err := service.Create(&req)

	assert.NoError(t, err, "should not error")
	assert.Equal(t, id, url.ID, fmt.Sprintf("id must be %d", id))
	assert.Equal(t, "/"+shortUrl, url.ShortUrl, fmt.Sprintf("shortUrl must be %s", shortUrl))
	assert.Equal(t, fullUrl, url.FullUrl, fmt.Sprintf("fullUrl must be %s", fullUrl))
}

func TestDeleteURLShouldNotFindInDB(t *testing.T) {
	repo := repositories.NewUrlRepositoryMock()
	service := services.NewUrlService(repo)

	id := uint64(1)

	err := service.Delete(id)
	assert.NoError(t, err, "should not error")

	_, err = service.GetById(id)
	assert.Error(t, err, "should error")
}
