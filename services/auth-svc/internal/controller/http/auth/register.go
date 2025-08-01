package auth

import (
	"github.com/gofiber/fiber/v2"
	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// Register handles user registration
// @Summary     Register new user
// @Description Register a new user account
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body model.RegisterRequest true "Registration data"
// @Success     201 {object} map[string]interface{} "User created successfully"
// @Failure     400 {object} map[string]string "Bad request"
// @Failure     409 {object} map[string]string "User already exists"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("Failed to parse request body", util.Error(err))
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get client info
	req.IPAddress = ctx.IP()
	req.UserAgent = ctx.Get("User-Agent")

	// Validate request
	if err := c.validateRegisterRequest(req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Register user
	user, err := c.authService.RegisterUser(ctx.Context(), req)
	if err != nil {
		c.logger.Error("Failed to register user", util.Error(err))

		// Handle specific errors
		switch err.Error() {
		case _const.CodeEmailExists.Message():
			return ctx.Status(_const.CodeEmailExists.HttpStatus()).JSON(fiber.Map{
				"code":    _const.CodeEmailExists.Code(),
				"message": _const.CodeEmailExists.Message(),
			})
		case _const.CodeUsernameExisted.Message():
			return ctx.Status(_const.CodeUsernameExisted.HttpStatus()).JSON(fiber.Map{
				"code":    _const.CodeUsernameExisted.Code(),
				"message": _const.CodeUsernameExisted.Message(),
			})
		default:
			return ctx.Status(500).JSON(fiber.Map{
				"code":    _const.CodeInternalError.Code(),
				"message": "Failed to register user",
			})
		}
	}

	return ctx.Status(201).JSON(fiber.Map{
		"code":    _const.CodeSuccess.Code(),
		"message": "User registered successfully",
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

func (c *AuthController) validateRegisterRequest(req model.RegisterRequest) error {
	if req.Email == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	if req.Username == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	if req.Password == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	if req.FirstName == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	if req.LastName == "" {
		return util.NewError(_const.CodeBadRequest.Message())
	}
	return nil
}
