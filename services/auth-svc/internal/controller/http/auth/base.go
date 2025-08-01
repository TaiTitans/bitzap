package auth

import (
	"github.com/taititans/bitzap/auth-svc/internal/service"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// AuthController handles HTTP requests for authentication
type AuthController struct {
	authService service.AuthService
	logger      util.Logger
}

// NewAuthController creates a new auth controller
func NewAuthController(authService service.AuthService, logger util.Logger) AuthControllerInterface {
	return &AuthController{
		authService: authService,
		logger:      logger,
	}
}
