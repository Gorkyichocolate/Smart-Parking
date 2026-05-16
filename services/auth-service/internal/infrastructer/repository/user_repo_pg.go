// internal/infrastructer/repository/user_repo_pg.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/domain"
)

// UserRepositoryPG реализация UserRepository для PostgreSQL
type UserRepositoryPG struct {
	db *sql.DB
}

// NewUserRepositoryPG создает новый экземпляр UserRepositoryPG
func NewUserRepositoryPG(db *sql.DB) *UserRepositoryPG {
	return &UserRepositoryPG{db: db}
}

// Create создает нового пользователя
func (r *UserRepositoryPG) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		string(user.Role),
		now,
	)

	if err != nil {
		// Проверяем на нарушение уникальности email
		if isUniqueViolation(err) {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	return nil
}

// GetByEmail возвращает пользователя по email
func (r *UserRepositoryPG) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepositoryPG) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// isUniqueViolation проверяет, является ли ошибка нарушением уникальности
func isUniqueViolation(err error) bool {
	errStr := err.Error()
	return contains(errStr, "duplicate key") || contains(errStr, "unique") || contains(errStr, "UNIQUE constraint")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && hasSubstring(s, substr)
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
