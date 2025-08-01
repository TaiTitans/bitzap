package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/config"
	"github.com/taititans/bitzap/auth-svc/internal/model"
)

// EmailRepository defines the interface for email operations
type EmailRepository interface {
	// Send generic email
	SendEmail(ctx context.Context, data model.EmailData) error

	// Get email configuration
	GetEmailConfig() config.EmailConfig
}
