package repository

import (
	"context"
	"errors"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"gorm.io/gorm"
)

// userRoleRepository implements UserRoleRepository
type userRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository creates a new user role repository
func NewUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return &userRoleRepository{db: db}
}

// Create creates a new user role
func (r *userRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

// GetByID gets user role by ID
func (r *userRoleRepository) GetByID(ctx context.Context, id uint) (*entity.UserRole, error) {
	var userRole entity.UserRole
	err := r.db.WithContext(ctx).First(&userRole, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userRole, nil
}

// Update updates user role
func (r *userRoleRepository) Update(ctx context.Context, userRole *entity.UserRole) error {
	return r.db.WithContext(ctx).Save(userRole).Error
}

// Delete deletes user role
func (r *userRoleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.UserRole{}, id).Error
}

// GetByUserID gets user roles by user ID
func (r *userRoleRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.UserRole, error) {
	var userRoles []*entity.UserRole
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

// AddRoleToUser adds role to user
func (r *userRoleRepository) AddRoleToUser(ctx context.Context, userID uint, role string) error {
	// Note: This is a simplified implementation. In a real app, you'd need to get role_id from roles table
	userRole := &entity.UserRole{
		UserID: userID,
		RoleID: 1, // Default role ID - you should implement proper role lookup
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

// RemoveRoleFromUser removes role from user
func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, userID uint, role string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role = ?", userID, role).Delete(&entity.UserRole{}).Error
}

// HasRole checks if user has role
func (r *userRoleRepository) HasRole(ctx context.Context, userID uint, role string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserRole{}).Where("user_id = ? AND role = ?", userID, role).Count(&count).Error
	return count > 0, err
}

// List gets user roles with pagination
func (r *userRoleRepository) List(ctx context.Context, offset, limit int) ([]*entity.UserRole, error) {
	var userRoles []*entity.UserRole
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&userRoles).Error
	return userRoles, err
}
