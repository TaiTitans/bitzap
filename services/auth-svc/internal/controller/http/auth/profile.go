package auth

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// GetProfile gets user profile
// @Summary     Get user profile
// @Description Get user profile by user ID
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       user_id path int true "User ID"
// @Success     200 {object} map[string]interface{} "User profile"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     404 {object} map[string]string "User not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/profile/{user_id} [get]
func (c *AuthController) GetProfile(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid user ID",
		})
	}

	user, err := c.authService.GetUserProfile(ctx.Context(), uint(userID))
	if err != nil {
		c.logger.Error("Failed to get user profile", util.Error(err))
		return ctx.Status(404).JSON(fiber.Map{
			"code":    _const.CodeUserNotFound.Code(),
			"message": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "User profile retrieved successfully",
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"username":    user.Username,
			"firstname":   user.Firstname,
			"lastname":    user.Lastname,
			"phone":       user.Phone,
			"is_active":   user.IsActive,
			"is_verified": user.IsVerified,
			"created_at":  user.CreatedAt,
			"updated_at":  user.UpdatedAt,
		},
	})
}

// UpdateProfile updates user profile
// @Summary     Update user profile
// @Description Update user profile information
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       user_id path int true "User ID"
// @Param       request body model.UpdateProfileRequest true "Profile update data"
// @Success     200 {object} map[string]interface{} "Profile updated successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     404 {object} map[string]string "User not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/profile/{user_id} [put]
func (c *AuthController) UpdateProfile(ctx *fiber.Ctx) error {
	userIDStr := ctx.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid user ID",
		})
	}

	var req model.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": "Invalid request body",
		})
	}

	// Validate request
	if err := c.validateUpdateProfileRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"code":    _const.CodeBadRequest.Code(),
			"message": err.Error(),
		})
	}

	// Update profile
	user, err := c.authService.UpdateUserProfile(ctx.Context(), uint(userID), req)
	if err != nil {
		c.logger.Error("Failed to update user profile", util.Error(err))
		return ctx.Status(404).JSON(fiber.Map{
			"code":    _const.CodeUserNotFound.Code(),
			"message": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"username":    user.Username,
			"firstname":   user.Firstname,
			"lastname":    user.Lastname,
			"phone":       user.Phone,
			"is_active":   user.IsActive,
			"is_verified": user.IsVerified,
			"updated_at":  user.UpdatedAt,
		},
	})
}

func (c *AuthController) validateUpdateProfileRequest(req model.UpdateProfileRequest) error {
	if req.FirstName == "" && req.LastName == "" && req.Phone == "" {
		return util.NewError("At least one field must be provided")
	}
	return nil
}
