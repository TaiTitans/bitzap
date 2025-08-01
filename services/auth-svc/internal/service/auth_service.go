package service

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/logic"
	"github.com/taititans/bitzap/auth-svc/internal/model"
)

// AuthService defines the interface for authentication service
type AuthService interface {
	// Register new user
	RegisterUser(ctx context.Context, req model.RegisterRequest) (*entity.User, error)

	// Login user
	LoginUser(ctx context.Context, req model.LoginRequest) (*entity.User, error)

	// Get user profile
	GetUserProfile(ctx context.Context, userID uint) (*entity.User, error)

	// Update user profile
	UpdateUserProfile(ctx context.Context, userID uint, req model.UpdateProfileRequest) (*entity.User, error)

	// Change password
	ChangePassword(ctx context.Context, userID uint, req model.ChangePasswordRequest) error

	// Request password reset
	RequestPasswordReset(ctx context.Context, req model.PasswordResetRequest) error

	// Reset password with token
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error

	// Verify email
	VerifyEmail(ctx context.Context, token string) error
}

// authService implements AuthService
type authService struct {
	authLogic *logic.AuthLogic
}

// NewAuthService creates a new auth service
func NewAuthService(authLogic *logic.AuthLogic) AuthService {
	return &authService{
		authLogic: authLogic,
	}
}

// RegisterUser registers a new user
func (s *authService) RegisterUser(ctx context.Context, req model.RegisterRequest) (*entity.User, error) {
	return s.authLogic.RegisterUser(ctx, req)
}

// LoginUser authenticates a user
func (s *authService) LoginUser(ctx context.Context, req model.LoginRequest) (*entity.User, error) {
	return s.authLogic.LoginUser(ctx, req)
}

// GetUserProfile gets user profile with roles and permissions
func (s *authService) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	return s.authLogic.GetUserProfile(ctx, userID)
}

// UpdateUserProfile updates user profile
func (s *authService) UpdateUserProfile(ctx context.Context, userID uint, req model.UpdateProfileRequest) (*entity.User, error) {
	return s.authLogic.UpdateUserProfile(ctx, userID, req)
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, userID uint, req model.ChangePasswordRequest) error {
	return s.authLogic.ChangePassword(ctx, userID, req)
}

// RequestPasswordReset requests password reset
func (s *authService) RequestPasswordReset(ctx context.Context, req model.PasswordResetRequest) error {
	return s.authLogic.RequestPasswordReset(ctx, req)
}

// ResetPassword resets user password with token
func (s *authService) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	return s.authLogic.ResetPassword(ctx, req)
}

// VerifyEmail verifies user email
func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	return s.authLogic.VerifyEmail(ctx, token)
}
