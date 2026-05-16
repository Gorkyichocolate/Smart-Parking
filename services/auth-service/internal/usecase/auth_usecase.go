package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/domain"
	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/repository"
)

var (
	secretKey  = []byte("smart-parking-secret-key-2024")
	accessTTL  = 15 * time.Minute
	refreshTTL = 168 * time.Hour
	issuerName = "smart-parking"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Role     domain.UserRole
}

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	if input.Email == "" || input.Password == "" || input.FullName == "" {
		return nil, domain.ErrInvalidUserData
	}
	if len(input.Password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
		Role:         input.Role,
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := uc.generateToken(user, accessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.generateToken(user, refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

type ValidateTokenInput struct {
	Token string
}

type ValidateTokenOutput struct {
	Valid  bool
	UserID string
	Role   string
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, input ValidateTokenInput) (*ValidateTokenOutput, error) {
	token, err := jwt.Parse(input.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return &ValidateTokenOutput{Valid: false}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &ValidateTokenOutput{Valid: false}, nil
	}

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	return &ValidateTokenOutput{
		Valid:  true,
		UserID: userID,
		Role:   role,
	}, nil
}

func (uc *AuthUseCase) generateToken(user *domain.User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"role":      string(user.Role),
		"full_name": user.FullName,
		"iss":       issuerName,
		"iat":       now.Unix(),
		"exp":       now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
