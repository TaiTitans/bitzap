package email

import "github.com/gofiber/fiber/v2"

// EmailControllerInterface defines the interface for email controller
type EmailControllerInterface interface {
	SendEmailVerification(ctx *fiber.Ctx) error
	SendPasswordReset(ctx *fiber.Ctx) error
	SendWelcomeEmail(ctx *fiber.Ctx) error
	SendCustomEmail(ctx *fiber.Ctx) error
}
