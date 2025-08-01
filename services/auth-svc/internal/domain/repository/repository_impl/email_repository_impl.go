package repository

import (
	"context"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/taititans/bitzap/auth-svc/internal/config"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// emailRepository implements repository.EmailRepository
type emailRepository struct {
	mailjetClient *mailjet.Client
	config        config.EmailConfig
	logger        util.Logger
}

// NewEmailRepository creates a new email repository
func NewEmailRepository(config config.EmailConfig, logger util.Logger) repository.EmailRepository {
	client := mailjet.NewMailjetClient(config.MailjetAPIKey, config.MailjetSecretKey)

	return &emailRepository{
		mailjetClient: client,
		config:        config,
		logger:        logger,
	}
}

// SendEmail sends email using Mailjet
func (r *emailRepository) SendEmail(ctx context.Context, data model.EmailData) error {
	r.logger.Info("Sending email via repository",
		util.String("to_email", data.ToEmail),
		util.String("subject", data.Subject),
	)

	// TODO: Implement actual Mailjet sending
	// For now, we'll just simulate success
	r.logger.Info("Email sent successfully (simulated)",
		util.String("to_email", data.ToEmail),
		util.String("subject", data.Subject),
	)

	return nil
}

// GetEmailConfig returns email configuration
func (r *emailRepository) GetEmailConfig() config.EmailConfig {
	return r.config
}
