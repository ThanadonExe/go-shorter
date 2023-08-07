package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/thanadonexe/go-shorter/internal/core/domain"
	"github.com/thanadonexe/go-shorter/internal/core/ports"

	_ "github.com/lib/pq"
)

const (
	UrlCacheDuration time.Duration = 24 * time.Hour
	UrlCachePrefix   string        = "url_id"
)

type UrlRepository struct {
	db    *sqlx.DB
	cache ports.CacheRepository
}

func NewUrlRepository(db *sqlx.DB, cache ports.CacheRepository) *UrlRepository {
	return &UrlRepository{
		db:    db,
		cache: cache,
	}
}

func (r *UrlRepository) NextUrlID() (uint64, error) {
	var id uint64
	query := `SELECT nextval('urls_id_seq');`
	err := r.db.QueryRowx(query).Scan(&id)

	return id, err
}

func (r *UrlRepository) Create(req *domain.UrlCreateRequest) (*domain.Url, error) {
	query := `INSERT INTO urls (id, short_url, full_url) VALUES ($1, $2, $3) RETURNING id, short_url, full_url, created_at, updated_at;`
	rows, err := r.db.Queryx(query, req.ID, req.ShortUrl, req.FullUrl)
	if err != nil {
		return nil, err
	}

	var url domain.Url
	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&url)

		if err != nil {
			return nil, err
		}
	}

	return &url, err
}

func (r *UrlRepository) GetById(id uint64) (*domain.Url, error) {
	query := `SELECT id, short_url, full_url, created_at, updated_at FROM urls WHERE id = :id AND deleted_at IS NULL;`
	rows, err := r.db.NamedQuery(query, map[string]interface{}{"id": id})

	if err != nil {
		return nil, err
	}

	var url domain.Url
	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&url)

		if err != nil {
			return nil, err
		}
	}

	return &url, err
}

func (r *UrlRepository) Delete(id uint64) error {
	query := `UPDATE urls SET deleted_at = NOW() WHERE id = :id;`

	_, err := r.db.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *UrlRepository) GetAll(offset, limit int) ([]*domain.Url, error) {
	query := `SELECT id, short_url, full_url, created_at, updated_at FROM urls WHERE deleted_at IS NULL OFFSET :offset LIMIT :limit;`

	rows, err := r.db.NamedQuery(query, map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	})

	if err != nil {
		return nil, err
	}

	var urls []*domain.Url
	err = sql.ErrNoRows
	for rows.Next() {
		var u domain.Url
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &u)
	}

	return urls, err
}

func (r *UrlRepository) GetByCode(code string) (*domain.Url, error) {
	var url domain.Url
	cacheKey := fmt.Sprintf("%s_%s", UrlCachePrefix, code)
	err := r.cache.Get(cacheKey, &url)
	if err == nil {
		return &url, nil
	}

	query := `SELECT id, short_url, full_url, created_at, updated_at FROM urls WHERE short_url = :code AND deleted_at IS NULL;`
	rows, err := r.db.NamedQuery(query, map[string]interface{}{"code": code})
	if err != nil {
		return nil, err
	}

	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&url)
		if err != nil {
			return nil, err
		}

		_ = r.cache.Set(cacheKey, &url, UrlCacheDuration)
	}

	return &url, err
}

func (r *UrlRepository) Count() int {
	var count int

	query := `SELECT COUNT(id) FROM urls WHERE deleted_at IS NULL;`
	err := r.db.Get(&count, query)
	if err != nil {
		return 0
	}

	return count
}
