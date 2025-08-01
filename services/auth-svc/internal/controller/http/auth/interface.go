package auth

import "github.com/gofiber/fiber/v2"

// AuthControllerInterface defines the interface for auth controller
type AuthControllerInterface interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	GetProfile(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
	RequestPasswordReset(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	VerifyEmail(ctx *fiber.Ctx) error
}
