// internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/domain"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
