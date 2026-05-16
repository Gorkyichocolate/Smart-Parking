package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth"
	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/domain"
	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/usecase"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	input := usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     domain.UserRole(req.Role),
	}

	user, err := h.authUseCase.Register(ctx, input)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &authpb.RegisterResponse{
		UserId:  user.ID,
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authUseCase.Login(ctx, input)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}

	log.Printf("User logged in: %s", req.Email)

	return &authpb.LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		ExpiresAt:    timestamppb.New(time.Now().Add(15 * time.Minute)),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	input := usecase.ValidateTokenInput{
		Token: req.Token,
	}

	output, err := h.authUseCase.ValidateToken(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate token: %v", err)
	}

	return &authpb.ValidateTokenResponse{
		Valid:  output.Valid,
		UserId: output.UserID,
		Role:   output.Role,
	}, nil
}
