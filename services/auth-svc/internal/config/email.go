package config

// EmailConfig holds email service configuration
type EmailConfig struct {
	MailjetAPIKey    string `yaml:"mailjet_api_key" env:"MAILJET_API_KEY"`
	MailjetSecretKey string `yaml:"mailjet_secret_key" env:"MAILJET_SECRET_KEY"`
	FromEmail        string `yaml:"from_email" env:"FROM_EMAIL"`
	FromName         string `yaml:"from_name" env:"FROM_NAME"`
	AppURL           string `yaml:"app_url" env:"APP_URL"`
}

// DefaultEmailConfig returns default email configuration
func DefaultEmailConfig() EmailConfig {
	return EmailConfig{
		FromEmail: "bitzapofficial@gmail.com",
		FromName:  "Bitzap Auth Official",
		AppURL:    "http://localhost:8080",
	}
}
