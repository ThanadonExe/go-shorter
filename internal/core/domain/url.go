package domain

import (
	"fmt"
	"time"

	"github.com/thanadonexe/go-shorter/internal/config"
)

type Url struct {
	ID        uint64    `json:"id" db:"id"`
	ShortUrl  string    `json:"short_url" db:"short_url"`
	FullUrl   string    `json:"full_url" db:"full_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}

type UrlResponse struct {
	ID        uint64    `json:"id" db:"id"`
	ShortUrl  string    `json:"short_url" db:"short_url"`
	FullUrl   string    `json:"full_url" db:"full_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UrlListResponse struct {
	Data      []*UrlResponse `json:"data"`
	Total     int            `json:"total"`
	Limit     int            `json:"limit"`
	Page      int            `json:"current_page"`
	TotalPage int            `json:"total_page"`
}

type UrlCreateRequest struct {
	ID       uint64 `json:"id,omitempty" db:"id"`
	ShortUrl string `json:"short_url,omitempty" db:"short_url"`
	FullUrl  string `json:"full_url" db:"full_url" validate:"required"`
}

func (u *Url) ToUrlResponse() *UrlResponse {
	shortUrl := fmt.Sprintf("%s/%s", config.AppConfig.Domain, u.ShortUrl)

	return &UrlResponse{
		ID:        u.ID,
		ShortUrl:  shortUrl,
		FullUrl:   u.FullUrl,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
