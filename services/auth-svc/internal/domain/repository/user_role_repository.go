package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
)

// UserRoleRepository defines the interface for user role data access
type UserRoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, userRole *entity.UserRole) error
	GetByID(ctx context.Context, id uint) (*entity.UserRole, error)
	Update(ctx context.Context, userRole *entity.UserRole) error
	Delete(ctx context.Context, id uint) error

	// User role operations
	GetByUserID(ctx context.Context, userID uint) ([]*entity.UserRole, error)
	AddRoleToUser(ctx context.Context, userID uint, role string) error
	RemoveRoleFromUser(ctx context.Context, userID uint, role string) error
	HasRole(ctx context.Context, userID uint, role string) (bool, error)

	// List operations
	List(ctx context.Context, offset, limit int) ([]*entity.UserRole, error)
}
