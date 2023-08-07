package domain

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserResponse struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserListResponse struct {
	Data      []*UserResponse `json:"data"`
	Total     int             `json:"total"`
	Limit     int             `json:"limit"`
	Page      int             `json:"current_page"`
	TotalPage int             `json:"total_page"`
}

type UserCreateRequest struct {
	Email           string `json:"email" db:"email" validate:"required,email"`
	Password        string `json:"password" db:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" db:"first_name" validate:"required"`
	LastName        string `json:"last_name" db:"last_name" validate:"required"`
}

type UserUpdateRequest struct {
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func (u User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) ApplyUpdate(req *UserUpdateRequest) {
	if req.Password != "" {
		u.Password = req.Password
	}

	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}

	if req.LastName != "" {
		u.LastName = req.LastName
	}

	u.UpdatedAt = time.Now().UTC()
}
