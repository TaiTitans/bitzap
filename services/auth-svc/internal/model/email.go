package model

// EmailVerificationRequest represents email verification request
type EmailVerificationRequest struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// EmailTemplate represents email template data
type EmailTemplate struct {
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Data    map[string]string `json:"data"`
}

// EmailData represents email sending data
type EmailData struct {
	ToEmail   string            `json:"to_email"`
	ToName    string            `json:"to_name"`
	Subject   string            `json:"subject"`
	HTMLBody  string            `json:"html_body"`
	TextBody  string            `json:"text_body"`
	Variables map[string]string `json:"variables"`
}
