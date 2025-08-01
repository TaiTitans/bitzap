package auth

import (
	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// VerifyEmail verifies user email
// @Summary Verify user email
// @Description Verify user email with token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/verify-email [get]
func (c *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Token is required",
		})
	}

	// Verify email
	err := c.authService.VerifyEmail(ctx.Context(), token)
	if err != nil {
		c.logger.Error("Failed to verify email", util.Error(err))
		return ctx.Status(500).JSON(fiber.Map{
			"code":    _const.CodeInternalError.Code(),
			"message": "Failed to verify email",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Email verified successfully",
	})
}
