package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
)

// UserPermissionRepository defines the interface for user permission data access
type UserPermissionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, permission *entity.UserPermission) error
	GetByID(ctx context.Context, id uint) (*entity.UserPermission, error)
	Update(ctx context.Context, permission *entity.UserPermission) error
	Delete(ctx context.Context, id uint) error

	// User permission operations
	GetByUserID(ctx context.Context, userID uint) ([]*entity.UserPermission, error)
	AddPermissionToUser(ctx context.Context, userID uint, resource, action string) error
	RemovePermissionFromUser(ctx context.Context, userID uint, resource, action string) error
	HasPermission(ctx context.Context, userID uint, resource, action string) (bool, error)

	// List operations
	List(ctx context.Context, offset, limit int) ([]*entity.UserPermission, error)
}
