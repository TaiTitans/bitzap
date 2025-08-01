package email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// SendEmailVerification handles email verification request
// @Summary Send email verification
// @Description Send verification email to user
// @Tags email
// @Accept json
// @Produce json
// @Param request body model.EmailVerificationRequest true "Email verification request"
// @Success 200 {object} map[string]string "Email sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /email/verify [post]
func (c *EmailController) SendEmailVerification(ctx *fiber.Ctx) error {
	var req model.EmailVerificationRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse email verification request", util.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and user_id are required",
		})
	}

	if err := c.emailService.SendEmailVerification(ctx.Context(), req); err != nil {
		c.logger.Error("Failed to send email verification", util.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send verification email",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Verification email sent successfully",
	})
}
