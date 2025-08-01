package auth

import (
	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// RequestPasswordReset requests password reset
// @Summary     Request password reset
// @Description Send password reset email to user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.PasswordResetRequest true "Password reset request"
// @Success     200 {object} map[string]interface{} "Password reset email sent"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/forgot-password [post]
func (c *AuthController) RequestPasswordReset(ctx *fiber.Ctx) error {
	// Parse request
	var req model.PasswordResetRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := c.validatePasswordResetRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": err.Error(),
		})
	}

	// Request password reset
	err := c.authService.RequestPasswordReset(ctx.Context(), req)
	if err != nil {
		c.logger.Error("Failed to request password reset", util.Error(err))
		return ctx.Status(500).JSON(fiber.Map{
			"code":    _const.CodeInternalError.Code(),
			"message": "Failed to request password reset",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Password reset email sent successfully",
	})
}

func (c *AuthController) validatePasswordResetRequest(req model.PasswordResetRequest) error {
	if req.Email == "" {
		return util.NewError("Email is required")
	}
	return nil
}
