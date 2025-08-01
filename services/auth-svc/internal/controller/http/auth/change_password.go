package auth

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// ChangePassword changes user password
// @Summary     Change user password
// @Description Change user password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       user_id path int true "User ID"
// @Param       request body model.ChangePasswordRequest true "Password change data"
// @Success     200 {object} map[string]interface{} "Password changed successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     404 {object} map[string]string "User not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/password/{user_id} [put]
func (c *AuthController) ChangePassword(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid user ID",
		})
	}

	var req model.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := c.validateChangePasswordRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": err.Error(),
		})
	}

	// Change password
	err = c.authService.ChangePassword(ctx.Context(), uint(userID), req)
	if err != nil {
		c.logger.Error("Failed to change password", util.Error(err))
		return ctx.Status(404).JSON(fiber.Map{
			"code":    _const.CodeUserNotFound.Code(),
			"message": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Password changed successfully",
	})
}

func (c *AuthController) validateChangePasswordRequest(req model.ChangePasswordRequest) error {
	if req.OldPassword == "" {
		return util.NewError("Old password is required")
	}
	if req.NewPassword == "" {
		return util.NewError("New password is required")
	}
	return nil
}
