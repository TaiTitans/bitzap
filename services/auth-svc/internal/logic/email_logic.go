package logic

import (
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// EmailLogic contains email business logic
type EmailLogic struct {
	emailRepo repository.EmailRepository
	redisRepo repository.RedisRepository
	logger    util.Logger
}

// NewEmailLogic creates a new email logic instance
func NewEmailLogic(emailRepo repository.EmailRepository, redisRepo repository.RedisRepository, logger util.Logger) *EmailLogic {
	return &EmailLogic{
		emailRepo: emailRepo,
		redisRepo: redisRepo,
		logger:    logger,
	}
}

// SendEmailVerification sends email verification
func (s *EmailLogic) SendEmailVerification(ctx context.Context, req model.EmailVerificationRequest) error {
	s.logger.Info("Sending email verification",
		util.String("email", req.Email),
		util.Int("user_id", int(req.UserID)),
	)

	// Get config from repository
	config := s.emailRepo.GetEmailConfig()

	// Generate UUID for verification token
	verificationToken := uuid.New().String()
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", config.AppURL, verificationToken)

	// Store token in Redis with expiration (24 hours)
	key := fmt.Sprintf("email_verification:%s", verificationToken)
	value := fmt.Sprintf("%d", req.UserID)
	expiration := 24 * time.Hour

	if err := s.redisRepo.Set(ctx, key, value, expiration); err != nil {
		s.logger.Error("Failed to store verification token in Redis", util.Error(err))
		return err
	}

	// Email template data
	data := map[string]string{
		"Name":            "User",
		"VerificationURL": verificationURL,
		"AppName":         "Bitzap",
		"SupportEmail":    "support@bitzap.com",
	}

	emailData := model.EmailData{
		ToEmail:   req.Email,
		ToName:    "User",
		Subject:   "Verify Your Email - Bitzap",
		HTMLBody:  s.generateEmailVerificationHTML(data),
		TextBody:  s.generateEmailVerificationText(data),
		Variables: data,
	}

	return s.emailRepo.SendEmail(ctx, emailData)
}

// SendPasswordReset sends password reset email
func (s *EmailLogic) SendPasswordReset(ctx context.Context, req model.PasswordResetRequest) error {
	s.logger.Info("Sending password reset email",
		util.String("email", req.Email),
	)

	// Get config from repository
	config := s.emailRepo.GetEmailConfig()

	// Generate UUID for reset token
	resetToken := uuid.New().String()
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", config.AppURL, resetToken)

	// Store token in Redis with expiration (1 hour)
	key := fmt.Sprintf("password_reset:%s", resetToken)
	value := req.Email
	expiration := 1 * time.Hour

	if err := s.redisRepo.Set(ctx, key, value, expiration); err != nil {
		s.logger.Error("Failed to store reset token in Redis", util.Error(err))
		return err
	}

	// Email template data
	data := map[string]string{
		"Name":         "User",
		"ResetURL":     resetURL,
		"AppName":      "Bitzap",
		"SupportEmail": "support@bitzap.com",
	}

	emailData := model.EmailData{
		ToEmail:   req.Email,
		ToName:    "User",
		Subject:   "Reset Your Password - Bitzap",
		HTMLBody:  s.generatePasswordResetHTML(data),
		TextBody:  s.generatePasswordResetText(data),
		Variables: data,
	}

	return s.emailRepo.SendEmail(ctx, emailData)
}

// SendWelcomeEmail sends welcome email
func (s *EmailLogic) SendWelcomeEmail(ctx context.Context, email, name string) error {
	s.logger.Info("Sending welcome email",
		util.String("email", email),
		util.String("name", name),
	)

	// Get config from repository
	config := s.emailRepo.GetEmailConfig()

	data := map[string]string{
		"Name":         name,
		"AppName":      "Bitzap",
		"LoginURL":     fmt.Sprintf("%s/auth/login", config.AppURL),
		"SupportEmail": "support@bitzap.com",
	}

	emailData := model.EmailData{
		ToEmail:   email,
		ToName:    name,
		Subject:   "Welcome to Bitzap!",
		HTMLBody:  s.generateWelcomeHTML(data),
		TextBody:  s.generateWelcomeText(data),
		Variables: data,
	}

	return s.emailRepo.SendEmail(ctx, emailData)
}

// SendEmail sends generic email using Mailjet
func (s *EmailLogic) SendEmail(ctx context.Context, data model.EmailData) error {
	s.logger.Info("Sending email",
		util.String("to_email", data.ToEmail),
		util.String("subject", data.Subject),
	)

	return s.emailRepo.SendEmail(ctx, data)
}

// VerifyEmailToken verifies email verification token from Redis
func (s *EmailLogic) VerifyEmailToken(ctx context.Context, token string) (uint, error) {
	key := fmt.Sprintf("email_verification:%s", token)

	// Get user ID from Redis
	userIDStr, err := s.redisRepo.Get(ctx, key)
	if err != nil {
		s.logger.Error("Failed to get verification token from Redis", util.Error(err))
		return 0, err
	}

	if userIDStr == "" {
		return 0, fmt.Errorf("invalid or expired verification token")
	}

	// Parse user ID
	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		s.logger.Error("Failed to parse user ID from token", util.Error(err))
		return 0, err
	}

	// Delete token from Redis after successful verification
	if err := s.redisRepo.Del(ctx, key); err != nil {
		s.logger.Error("Failed to delete verification token from Redis", util.Error(err))
		// Don't fail verification if deletion fails
	}

	return userID, nil
}

// VerifyPasswordResetToken verifies password reset token from Redis
func (s *EmailLogic) VerifyPasswordResetToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("password_reset:%s", token)

	// Get email from Redis
	email, err := s.redisRepo.Get(ctx, key)
	if err != nil {
		s.logger.Error("Failed to get reset token from Redis", util.Error(err))
		return "", err
	}

	if email == "" {
		return "", fmt.Errorf("invalid or expired reset token")
	}

	// Delete token from Redis after successful verification
	if err := s.redisRepo.Del(ctx, key); err != nil {
		s.logger.Error("Failed to delete reset token from Redis", util.Error(err))
		// Don't fail verification if deletion fails
	}

	return email, nil
}

// generateEmailVerificationHTML generates HTML for email verification
func (s *EmailLogic) generateEmailVerificationHTML(data map[string]string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.AppName}}</h1>
        </div>
        <div class="content">
            <h2>Verify Your Email Address</h2>
            <p>Hello {{.Name}},</p>
            <p>Thank you for signing up for {{.AppName}}! Please verify your email address by clicking the button below:</p>
            <p style="text-align: center;">
                <a href="{{.VerificationURL}}" class="button">Verify Email</a>
            </p>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p>{{.VerificationURL}}</p>
            <p>This link will expire in 24 hours.</p>
            <p>If you didn't create an account with {{.AppName}}, you can safely ignore this email.</p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
        </div>
    </div>
</body>
</html>`

	return s.renderTemplate(tmpl, data)
}

// generateEmailVerificationText generates text for email verification
func (s *EmailLogic) generateEmailVerificationText(data map[string]string) string {
	tmpl := `Verify Your Email Address

Hello {{.Name}},

Thank you for signing up for {{.AppName}}! Please verify your email address by visiting this link:

{{.VerificationURL}}

This link will expire in 24 hours.

If you didn't create an account with {{.AppName}}, you can safely ignore this email.

Need help? Contact us at {{.SupportEmail}}`

	return s.renderTemplate(tmpl, data)
}

// generatePasswordResetHTML generates HTML for password reset
func (s *EmailLogic) generatePasswordResetHTML(data map[string]string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #dc3545; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .button { display: inline-block; padding: 12px 24px; background: #dc3545; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.AppName}}</h1>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>Hello {{.Name}},</p>
            <p>We received a request to reset your password. Click the button below to create a new password:</p>
            <p style="text-align: center;">
                <a href="{{.ResetURL}}" class="button">Reset Password</a>
            </p>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p>{{.ResetURL}}</p>
            <p>This link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, you can safely ignore this email.</p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
        </div>
    </div>
</body>
</html>`

	return s.renderTemplate(tmpl, data)
}

// generatePasswordResetText generates text for password reset
func (s *EmailLogic) generatePasswordResetText(data map[string]string) string {
	tmpl := `Reset Your Password

Hello {{.Name}},

We received a request to reset your password. Visit this link to create a new password:

{{.ResetURL}}

This link will expire in 1 hour.

If you didn't request a password reset, you can safely ignore this email.

Need help? Contact us at {{.SupportEmail}}`

	return s.renderTemplate(tmpl, data)
}

// generateWelcomeHTML generates HTML for welcome email
func (s *EmailLogic) generateWelcomeHTML(data map[string]string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to {{.AppName}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f8f9fa; }
        .button { display: inline-block; padding: 12px 24px; background: #28a745; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to {{.AppName}}</h1>
        </div>
        <div class="content">
            <h2>Welcome aboard, {{.Name}}!</h2>
            <p>Thank you for joining {{.AppName}}! We're excited to have you as part of our community.</p>
            <p>You can now log in to your account and start exploring our features:</p>
            <p style="text-align: center;">
                <a href="{{.LoginURL}}" class="button">Login to Your Account</a>
            </p>
            <p>If you have any questions or need assistance, don't hesitate to reach out to our support team.</p>
        </div>
        <div class="footer">
            <p>Need help? Contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
        </div>
    </div>
</body>
</html>`

	return s.renderTemplate(tmpl, data)
}

// generateWelcomeText generates text for welcome email
func (s *EmailLogic) generateWelcomeText(data map[string]string) string {
	tmpl := `Welcome to {{.AppName}}

Welcome aboard, {{.Name}}!

Thank you for joining {{.AppName}}! We're excited to have you as part of our community.

You can now log in to your account and start exploring our features:

{{.LoginURL}}

If you have any questions or need assistance, don't hesitate to reach out to our support team.

Need help? Contact us at {{.SupportEmail}}`

	return s.renderTemplate(tmpl, data)
}

// renderTemplate renders template with data
func (s *EmailLogic) renderTemplate(tmpl string, data map[string]string) string {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		s.logger.Error("Failed to parse template", util.Error(err))
		return ""
	}

	var buf strings.Builder
	err = t.Execute(&buf, data)
	if err != nil {
		s.logger.Error("Failed to execute template", util.Error(err))
		return ""
	}

	return buf.String()
}
