package logic

import (
	"context"

	_const "github.com/taititans/bitzap/auth-svc/internal/const"
	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"github.com/taititans/bitzap/auth-svc/internal/model"
	"github.com/taititans/bitzap/auth-svc/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// EmailServiceInterface defines the interface for email service
type EmailServiceInterface interface {
	SendEmailVerification(ctx context.Context, req model.EmailVerificationRequest) error
	SendPasswordReset(ctx context.Context, req model.PasswordResetRequest) error
	SendWelcomeEmail(ctx context.Context, email, name string) error
	SendEmail(ctx context.Context, data model.EmailData) error
	VerifyEmailToken(ctx context.Context, token string) (uint, error)
	VerifyPasswordResetToken(ctx context.Context, token string) (string, error)
}

// AuthLogic contains business logic for authentication
type AuthLogic struct {
	userRepo           repository.UserRepository
	userRoleRepo       repository.UserRoleRepository
	userPermissionRepo repository.UserPermissionRepository
	userActivityRepo   repository.UserActivityLogRepository
	emailService       EmailServiceInterface
	logger             util.Logger
}

// NewAuthLogic creates new AuthLogic instance
func NewAuthLogic(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	userPermissionRepo repository.UserPermissionRepository,
	userActivityRepo repository.UserActivityLogRepository,
	emailService EmailServiceInterface,
	logger util.Logger,
) *AuthLogic {
	return &AuthLogic{
		userRepo:           userRepo,
		userRoleRepo:       userRoleRepo,
		userPermissionRepo: userPermissionRepo,
		userActivityRepo:   userActivityRepo,
		emailService:       emailService,
		logger:             logger,
	}
}

// RegisterUser registers a new user
func (l *AuthLogic) RegisterUser(ctx context.Context, req model.RegisterRequest) (*entity.User, error) {
	l.logger.Info("Registering new user",
		util.String("email", req.Email),
		util.String("username", req.Username),
	)

	// Check if email already exists
	existingUser, err := l.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		l.logger.Error("Failed to check existing email", util.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, util.NewError(_const.CodeEmailExists.Message())
	}

	// Check if username already exists
	existingUser, err = l.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		l.logger.Error("Failed to check existing username", util.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, util.NewError(_const.CodeUsernameExisted.Message())
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.logger.Error("Failed to hash password", util.Error(err))
		return nil, err
	}

	// Create user
	user := &entity.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Firstname:    req.FirstName,
		Lastname:     req.LastName,
		Phone:        req.Phone,
		IsActive:     true,
		IsVerified:   false,
	}

	if err := l.userRepo.Create(ctx, user); err != nil {
		l.logger.Error("Failed to create user", util.Error(err))
		return nil, err
	}

	// Send welcome email
	if err := l.emailService.SendWelcomeEmail(ctx, user.Email, user.Firstname); err != nil {
		l.logger.Error("Failed to send welcome email", util.Error(err))
		// Don't fail registration if email fails
	}

	// Send email verification
	if err := l.emailService.SendEmailVerification(ctx, model.EmailVerificationRequest{
		UserID: user.ID,
		Email:  user.Email,
	}); err != nil {
		l.logger.Error("Failed to send email verification", util.Error(err))
		// Don't fail registration if email fails
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, user.ID, "register", "user", req.IPAddress, req.UserAgent, nil)

	l.logger.Info("User registered successfully",
		util.Int("user_id", int(user.ID)),
		util.String("email", user.Email),
	)

	return user, nil
}

// LoginUser authenticates a user
func (l *AuthLogic) LoginUser(ctx context.Context, req model.LoginRequest) (*entity.User, error) {
	l.logger.Info("User login attempt",
		util.String("email", req.Email),
	)

	user, err := l.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		l.logger.Error("Failed to get user by email", util.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, util.NewError(_const.CodeWrongPassword.Message())
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		l.logger.Warn("Invalid password for user",
			util.String("email", req.Email),
		)
		return nil, util.NewError(_const.CodeWrongPassword.Message())
	}

	// Check if user is active
	if !user.IsActive {
		return nil, util.NewError(_const.CodeLockingAccount.Message())
	}

	// Update last login
	if err := l.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		l.logger.Error("Failed to update last login", util.Error(err))
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, user.ID, "login", "user", req.IPAddress, req.UserAgent, nil)

	l.logger.Info("User logged in successfully",
		util.Int("user_id", int(user.ID)),
		util.String("email", user.Email),
	)

	return user, nil
}

// GetUserProfile gets user profile with roles and permissions
func (l *AuthLogic) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	l.logger.Info("Getting user profile",
		util.Int("user_id", int(userID)),
	)

	user, err := l.userRepo.GetWithRoles(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to get user with roles", util.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, util.NewError(_const.CodeUserNotFound.Message())
	}

	// Get permissions
	permissions, err := l.userPermissionRepo.GetByUserID(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to get user permissions", util.Error(err))
		return nil, err
	}
	// Convert []*UserPermission to []UserPermission
	userPermissions := make([]entity.UserPermission, len(permissions))
	for i, p := range permissions {
		userPermissions[i] = *p
	}
	user.Permissions = userPermissions

	return user, nil
}

// UpdateUserProfile updates user profile
func (l *AuthLogic) UpdateUserProfile(ctx context.Context, userID uint, req model.UpdateProfileRequest) (*entity.User, error) {
	l.logger.Info("Updating user profile",
		util.Int("user_id", int(userID)),
	)

	user, err := l.userRepo.GetByID(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to get user", util.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, util.NewError(_const.CodeUserNotFound.Message())
	}

	// Update fields
	user.Firstname = req.FirstName
	user.Lastname = req.LastName
	user.Phone = req.Phone
	user.AvatarURL = req.AvatarURL

	if err := l.userRepo.Update(ctx, user); err != nil {
		l.logger.Error("Failed to update user", util.Error(err))
		return nil, err
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, userID, "update_profile", "user", req.IPAddress, req.UserAgent, nil)

	l.logger.Info("User profile updated successfully",
		util.Int("user_id", int(userID)),
	)

	return user, nil
}

// ChangePassword changes user password
func (l *AuthLogic) ChangePassword(ctx context.Context, userID uint, req model.ChangePasswordRequest) error {
	l.logger.Info("Changing user password",
		util.Int("user_id", int(userID)),
	)

	user, err := l.userRepo.GetByID(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to get user", util.Error(err))
		return err
	}
	if user == nil {
		return util.NewError(_const.CodeUserNotFound.Message())
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		l.logger.Warn("Invalid old password for user",
			util.Int("user_id", int(userID)),
		)
		return util.NewError(_const.CodeWrongOldPassword.Message())
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.logger.Error("Failed to hash new password", util.Error(err))
		return err
	}

	user.PasswordHash = string(hashedPassword)
	if err := l.userRepo.Update(ctx, user); err != nil {
		l.logger.Error("Failed to update password", util.Error(err))
		return err
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, userID, "change_password", "user", req.IPAddress, req.UserAgent, nil)

	l.logger.Info("Password changed successfully",
		util.Int("user_id", int(userID)),
	)

	return nil
}

// RequestPasswordReset requests password reset
func (l *AuthLogic) RequestPasswordReset(ctx context.Context, req model.PasswordResetRequest) error {
	l.logger.Info("Requesting password reset",
		util.String("email", req.Email),
	)

	// Check if user exists
	user, err := l.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		l.logger.Error("Failed to get user by email", util.Error(err))
		return err
	}
	if user == nil {
		// Don't reveal if user exists or not for security
		l.logger.Info("Password reset requested for non-existent email",
			util.String("email", req.Email),
		)
		return nil
	}

	// Send password reset email
	if err := l.emailService.SendPasswordReset(ctx, req); err != nil {
		l.logger.Error("Failed to send password reset email", util.Error(err))
		return err
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, user.ID, "password_reset_request", "user", "", "", nil)

	l.logger.Info("Password reset email sent successfully",
		util.String("email", req.Email),
	)

	return nil
}

// VerifyEmail verifies user email
func (l *AuthLogic) VerifyEmail(ctx context.Context, token string) error {
	l.logger.Info("Verifying email with token",
		util.String("token", token[:10]+"..."), // Log first 10 chars for security
	)

	// Verify token using email logic
	userID, err := l.emailService.VerifyEmailToken(ctx, token)
	if err != nil {
		l.logger.Error("Failed to verify email token", util.Error(err))
		return err
	}

	// Update user verification status in database
	user, err := l.userRepo.GetByID(ctx, userID)
	if err != nil {
		l.logger.Error("Failed to get user by ID", util.Error(err))
		return err
	}
	if user == nil {
		return util.NewError("User not found")
	}

	user.IsVerified = true
	if err := l.userRepo.Update(ctx, user); err != nil {
		l.logger.Error("Failed to update user verification status", util.Error(err))
		return err
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, userID, "email_verified", "user", "", "", nil)

	l.logger.Info("Email verified successfully",
		util.Int("user_id", int(userID)),
	)

	return nil
}

// ResetPassword resets user password with token
func (l *AuthLogic) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	l.logger.Info("Resetting password with token",
		util.String("token", req.Token[:10]+"..."), // Log first 10 chars for security
	)

	// Verify reset token using email logic
	email, err := l.emailService.VerifyPasswordResetToken(ctx, req.Token)
	if err != nil {
		l.logger.Error("Failed to verify reset token", util.Error(err))
		return err
	}

	// Get user by email
	user, err := l.userRepo.GetByEmail(ctx, email)
	if err != nil {
		l.logger.Error("Failed to get user by email", util.Error(err))
		return err
	}
	if user == nil {
		return util.NewError("User not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.logger.Error("Failed to hash password", util.Error(err))
		return err
	}

	// Update user password
	user.PasswordHash = string(hashedPassword)
	if err := l.userRepo.Update(ctx, user); err != nil {
		l.logger.Error("Failed to update user password", util.Error(err))
		return err
	}

	// Log activity
	l.userActivityRepo.LogActivity(ctx, user.ID, "password_reset", "user", req.IPAddress, req.UserAgent, nil)

	l.logger.Info("Password reset successfully",
		util.Int("user_id", int(user.ID)),
		util.String("email", email),
	)

	return nil
}
