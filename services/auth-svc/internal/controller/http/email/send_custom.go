package email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// SendCustomEmail handles custom email request
// @Summary Send custom email
// @Description Send custom email with provided data
// @Tags email
// @Accept json
// @Produce json
// @Param request body model.EmailData true "Custom email data"
// @Success 200 {object} map[string]string "Email sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /email/send [post]
func (c *EmailController) SendCustomEmail(ctx *fiber.Ctx) error {
	var emailData model.EmailData
	if err := ctx.BodyParser(&emailData); err != nil {
		c.logger.Error("Failed to parse custom email request", util.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if emailData.ToEmail == "" || emailData.Subject == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ToEmail and Subject are required",
		})
	}

	if err := c.emailService.SendEmail(ctx.Context(), emailData); err != nil {
		c.logger.Error("Failed to send custom email", util.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}
