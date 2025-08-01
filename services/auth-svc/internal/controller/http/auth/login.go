package auth

import (
	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// Login handles user login
// @Summary     Login user
// @Description Authenticate user and return access token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.LoginRequest true "Login credentials"
// @Success     200 {object} map[string]interface{} "Login successful"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     401 {object} map[string]string "Invalid credentials"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req model.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := c.validateLoginRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Login user
	user, err := c.authService.LoginUser(ctx.Context(), req)
	if err != nil {
		c.logger.Error("Failed to login user", util.Error(err))
		return ctx.Status(401).JSON(fiber.Map{
			"code":    _const.CodeWrongPassword.Code(),
			"message": "Invalid email or password",
		})
	}

	return ctx.JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "Login successful",
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
			"username":    user.Username,
			"firstname":   user.Firstname,
			"lastname":    user.Lastname,
			"is_active":   user.IsActive,
			"is_verified": user.IsVerified,
		},
	})
}

func (c *AuthController) validateLoginRequest(req model.LoginRequest) error {
	if req.Email == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	if req.Password == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	return nil
}
