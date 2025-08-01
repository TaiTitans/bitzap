package email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// SendWelcomeEmail handles welcome email request
// @Summary Send welcome email
// @Description Send welcome email to new user
// @Tags email
// @Accept json
// @Produce json
// @Param request body map[string]string true "Welcome email request"
// @Success 200 {object} map[string]string "Email sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /email/welcome [post]
func (c *EmailController) SendWelcomeEmail(ctx *fiber.Ctx) error {
	var req map[string]string
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse welcome email request", util.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	email, ok := req["email"]
	if !ok || email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	name, ok := req["name"]
	if !ok || name == "" {
		name = "User" // Default name
	}

	if err := c.emailService.SendWelcomeEmail(ctx.Context(), email, name); err != nil {
		c.logger.Error("Failed to send welcome email", util.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send welcome email",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Welcome email sent successfully",
	})
}
