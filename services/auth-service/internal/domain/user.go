// internal/domain/user.go
package domain

import (
	"time"
)

// UserRole определяет роль пользователя
type UserRole string

const (
	RoleDriver UserRole = "driver"
	RoleAdmin  UserRole = "admin"
)

// User представляет пользователя системы
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// Domain errors
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound       = &DomainError{Code: "USER_NOT_FOUND", Message: "user not found"}
	ErrEmailAlreadyExists = &DomainError{Code: "EMAIL_ALREADY_EXISTS", Message: "email already exists"}
	ErrInvalidCredentials = &DomainError{Code: "INVALID_CREDENTIALS", Message: "invalid email or password"}
	ErrInvalidUserData    = &DomainError{Code: "INVALID_USER_DATA", Message: "invalid user data"}
)
