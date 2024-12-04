package dto

import (
	"time"
)

type UserCreateDTO struct {
	Name     string `json:"name" validate:"min=3,max=255"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=3,max=20"`
}

type UserUpdateDTO struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=3,max=20"`
}

type UserResponseDTO struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone,omitempty"`
	Password  string     `json:"password,omitempty"`
	RoleID    int        `json:"role_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
