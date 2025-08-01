package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
)

type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error

	// List operations
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, offset, limit int) ([]*entity.User, error)

	// Relationship operations
	GetWithRoles(ctx context.Context, id uint) (*entity.User, error)
	GetWithPermissions(ctx context.Context, id uint) (*entity.User, error)
	GetWithActivityLogs(ctx context.Context, id uint, limit int) (*entity.User, error)

	// Authentication related
	UpdateLastLogin(ctx context.Context, id uint) error
	VerifyEmail(ctx context.Context, id uint) error
}
