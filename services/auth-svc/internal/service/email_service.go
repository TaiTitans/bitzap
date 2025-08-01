package service

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/taititans/bitzap/auth-svc/internal/config"
	repository_impl "github.com/taititans/bitzap/auth-svc/internal/domain/repository/repository_impl"
	"github.com/taititans/bitzap/auth-svc/internal/logic"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// EmailService defines the interface for email service
type EmailService interface {
	// Send email verification
	SendEmailVerification(ctx context.Context, req model.EmailVerificationRequest) error

	// Send password reset
	SendPasswordReset(ctx context.Context, req model.PasswordResetRequest) error

	// Send welcome email
	SendWelcomeEmail(ctx context.Context, email, name string) error

	// Send generic email
	SendEmail(ctx context.Context, data model.EmailData) error

	// Verify email token
	VerifyEmailToken(ctx context.Context, token string) (uint, error)

	// Verify password reset token
	VerifyPasswordResetToken(ctx context.Context, token string) (string, error)
}

// emailService implements EmailService
type emailService struct {
	emailLogic *logic.EmailLogic
	logger     util.Logger
}

// NewEmailService creates a new email service
func NewEmailService(config config.EmailConfig, redisClient *redis.Client, logger util.Logger) EmailService {
	// Create email repository
	emailRepo := repository_impl.NewEmailRepository(config, logger)

	// Create Redis repository
	redisRepo := repository_impl.NewRedisRepository(redisClient, logger)

	// Create email logic
	emailLogic := logic.NewEmailLogic(emailRepo, redisRepo, logger)

	return &emailService{
		emailLogic: emailLogic,
		logger:     logger,
	}
}

// SendEmailVerification sends email verification
func (s *emailService) SendEmailVerification(ctx context.Context, req model.EmailVerificationRequest) error {
	return s.emailLogic.SendEmailVerification(ctx, req)
}

// SendPasswordReset sends password reset email
func (s *emailService) SendPasswordReset(ctx context.Context, req model.PasswordResetRequest) error {
	return s.emailLogic.SendPasswordReset(ctx, req)
}

// SendWelcomeEmail sends welcome email
func (s *emailService) SendWelcomeEmail(ctx context.Context, email, name string) error {
	return s.emailLogic.SendWelcomeEmail(ctx, email, name)
}

// SendEmail sends generic email using Mailjet
func (s *emailService) SendEmail(ctx context.Context, data model.EmailData) error {
	return s.emailLogic.SendEmail(ctx, data)
}

// VerifyEmailToken verifies email verification token
func (s *emailService) VerifyEmailToken(ctx context.Context, token string) (uint, error) {
	return s.emailLogic.VerifyEmailToken(ctx, token)
}

// VerifyPasswordResetToken verifies password reset token
func (s *emailService) VerifyPasswordResetToken(ctx context.Context, token string) (string, error) {
	return s.emailLogic.VerifyPasswordResetToken(ctx, token)
}
