package email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// SendPasswordReset handles password reset request
// @Summary Send password reset email
// @Description Send password reset email to user
// @Tags email
// @Accept json
// @Produce json
// @Param request body model.PasswordResetRequest true "Password reset request"
// @Success 200 {object} map[string]string "Email sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /email/reset-password [post]
func (c *EmailController) SendPasswordReset(ctx *fiber.Ctx) error {
	var req model.PasswordResetRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse password reset request", util.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	if err := c.emailService.SendPasswordReset(ctx.Context(), req); err != nil {
		c.logger.Error("Failed to send password reset email", util.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send password reset email",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Password reset email sent successfully",
	})
}
