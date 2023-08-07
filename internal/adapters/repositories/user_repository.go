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
	UserCacheDuration time.Duration = 24 * time.Hour
	UserCachePrefix   string        = "user_id"
)

type UserRepository struct {
	db    *sqlx.DB
	cache ports.CacheRepository
}

func NewUserRepository(db *sqlx.DB, cache ports.CacheRepository) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: cache,
	}
}

func (r *UserRepository) Create(request *domain.UserCreateRequest) (*domain.User, error) {
	query := `INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id, email, password, first_name, last_name, created_at, updated_at;`
	rows, err := r.db.Queryx(query, request.Email, request.Password, request.FirstName, request.LastName)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&user)

		if err != nil {
			return nil, err
		}
	}

	return &user, err
}

func (r *UserRepository) GetById(id int) (*domain.User, error) {
	var user domain.User
	cacheKey := fmt.Sprintf("%s_%d", UserCachePrefix, id)
	err := r.cache.Get(cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	query := `SELECT id, email, password, first_name, last_name, created_at, updated_at FROM users WHERE id = :id AND deleted_at IS NULL;`
	rows, err := r.db.NamedQuery(query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&user)

		if err != nil {
			return nil, err
		}
		_ = r.cache.Set(cacheKey, &user, UserCacheDuration)
	}

	return &user, err
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET password = :password, first_name = :first_name, last_name = :last_name, updated_at = :updated_at WHERE id = :id;`
	_, err := r.db.NamedExec(query, &user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = :id;`
	_, err := r.db.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetAll(offset, limit int) ([]*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, created_at, updated_at FROM users WHERE deleted_at IS NULL OFFSET :offset LIMIT :limit;`
	rows, err := r.db.NamedQuery(query, map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	err = sql.ErrNoRows
	for rows.Next() {
		var user domain.User
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, err
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, created_at, updated_at FROM users WHERE email = :email AND deleted_at IS NULL;`
	rows, err := r.db.NamedQuery(query, map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = sql.ErrNoRows
	for rows.Next() {
		err = rows.StructScan(&user)

		if err != nil {
			return nil, err
		}
	}

	return &user, err
}

func (r *UserRepository) Count() int {
	var count int

	query := `SELECT COUNT(id) FROM users WHERE deleted_at IS NULL;`
	err := r.db.Get(&count, query)
	if err != nil {
		return 0
	}

	return count
}
