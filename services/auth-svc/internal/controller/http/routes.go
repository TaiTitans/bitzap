package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/controller/http/auth"
	"github.com/taititans/bitzap/auth-svc/internal/controller/http/email"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(app *fiber.App, authController auth.AuthControllerInterface, emailController email.EmailControllerInterface) {
	// Auth group
	authGroup := app.Group("/auth")

	// Registration and login
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	// Email functionality
	authGroup.Post("/forgot-password", authController.RequestPasswordReset)
	authGroup.Post("/reset-password", authController.ResetPassword)
	authGroup.Get("/verify-email", authController.VerifyEmail)

	// Profile management
	authGroup.Get("/profile/:user_id", authController.GetProfile)
	authGroup.Put("/profile/:user_id", authController.UpdateProfile)
	authGroup.Put("/password/:user_id", authController.ChangePassword)

	// Email routes
	emailGroup := app.Group("/email")
	emailGroup.Post("/verify", emailController.SendEmailVerification)
	emailGroup.Post("/reset-password", emailController.SendPasswordReset)
	emailGroup.Post("/welcome", emailController.SendWelcomeEmail)
	emailGroup.Post("/send", emailController.SendCustomEmail)
}
