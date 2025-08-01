package repository

import (
	"context"
	"errors"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"gorm.io/gorm"
)

// userPermissionRepository implements UserPermissionRepository
type userPermissionRepository struct {
	db *gorm.DB
}

// NewUserPermissionRepository creates a new user permission repository
func NewUserPermissionRepository(db *gorm.DB) repository.UserPermissionRepository {
	return &userPermissionRepository{db: db}
}

// Create creates a new user permission
func (r *userPermissionRepository) Create(ctx context.Context, permission *entity.UserPermission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetByID gets user permission by ID
func (r *userPermissionRepository) GetByID(ctx context.Context, id uint) (*entity.UserPermission, error) {
	var permission entity.UserPermission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// Update updates user permission
func (r *userPermissionRepository) Update(ctx context.Context, permission *entity.UserPermission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// Delete deletes user permission
func (r *userPermissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.UserPermission{}, id).Error
}

// GetByUserID gets user permissions by user ID
func (r *userPermissionRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.UserPermission, error) {
	var permissions []*entity.UserPermission
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&permissions).Error
	return permissions, err
}

// AddPermissionToUser adds permission to user
func (r *userPermissionRepository) AddPermissionToUser(ctx context.Context, userID uint, resource, action string) error {
	permission := &entity.UserPermission{
		UserID:   userID,
		Resource: resource,
		Action:   action,
	}
	return r.db.WithContext(ctx).Create(permission).Error
}

// RemovePermissionFromUser removes permission from user
func (r *userPermissionRepository) RemovePermissionFromUser(ctx context.Context, userID uint, resource, action string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND resource = ? AND action = ?", userID, resource, action).Delete(&entity.UserPermission{}).Error
}

// HasPermission checks if user has permission
func (r *userPermissionRepository) HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserPermission{}).Where("user_id = ? AND resource = ? AND action = ?", userID, resource, action).Count(&count).Error
	return count > 0, err
}

// List gets user permissions with pagination
func (r *userPermissionRepository) List(ctx context.Context, offset, limit int) ([]*entity.UserPermission, error) {
	var permissions []*entity.UserPermission
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&permissions).Error
	return permissions, err
}
