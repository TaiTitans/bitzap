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

	// Check if Mailjet credentials are configured
	if r.config.MailjetAPIKey == "" || r.config.MailjetSecretKey == "" {
		r.logger.Warn("Mailjet credentials not configured, email sending disabled",
			util.String("api_key_set", r.getBoolString(r.config.MailjetAPIKey != "")),
			util.String("secret_key_set", r.getBoolString(r.config.MailjetSecretKey != "")),
		)
		r.logger.Info("Email would be sent (credentials missing)",
			util.String("to_email", data.ToEmail),
			util.String("subject", data.Subject),
		)
		return nil
	}

	// Create Mailjet message
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: r.config.FromEmail,
				Name:  r.config.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: data.ToEmail,
					Name:  data.ToName,
				},
			},
			Subject:  data.Subject,
			TextPart: data.TextBody,
			HTMLPart: data.HTMLBody,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send email
	res, err := r.mailjetClient.SendMailV31(&messages)
	if err != nil {
		r.logger.Error("Failed to send email via Mailjet",
			util.Error(err),
			util.String("to_email", data.ToEmail),
		)
		return err
	}

	r.logger.Info("Email sent successfully via Mailjet",
		util.String("to_email", data.ToEmail),
		util.String("subject", data.Subject),
		util.String("status", res.ResultsV31[0].Status),
	)

	return nil
}

// getBoolString converts boolean to string for logging
func (r *emailRepository) getBoolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// GetEmailConfig returns email configuration
func (r *emailRepository) GetEmailConfig() config.EmailConfig {
	return r.config
}
