package email

import (
	"github.com/taititans/bitzap/auth-svc/internal/service"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// EmailController handles email-related HTTP requests
type EmailController struct {
	emailService service.EmailService
	logger       util.Logger
}

// NewEmailController creates a new email controller
func NewEmailController(emailService service.EmailService, logger util.Logger) EmailControllerInterface {
	return &EmailController{
		emailService: emailService,
		logger:       logger,
	}
}
