package auth

import (
	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// ResetPassword handles password reset with token
// @Summary     Reset password with token
// @Description Reset user password using reset token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.ResetPasswordRequest true "Password reset data"
// @Success     200 {object} map[string]interface{} "Password reset successful"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	var req model.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := c.validateResetPasswordRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": err.Error(),
		})
	}

	// Reset password
	err := c.authService.ResetPassword(ctx.Context(), req)
	if err != nil {
		c.logger.Error("Failed to reset password", util.Error(err))
		return ctx.Status(500).JSON(fiber.Map{
			"code":    _const.CodeInternalError.Code(),
			"message": "Failed to reset password",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Password reset successfully",
	})
}

func (c *AuthController) validateResetPasswordRequest(req model.ResetPasswordRequest) error {
	if req.Token == "" {
		return util.NewError("Token is required")
	}
	if req.NewPassword == "" {
		return util.NewError("New password is required")
	}
	if req.ConfirmPassword == "" {
		return util.NewError("Confirm password is required")
	}
	if req.NewPassword != req.ConfirmPassword {
		return util.NewError("New password and confirm password must match")
	}
	return nil
}
