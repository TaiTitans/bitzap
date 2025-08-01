package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
)

// UserActivityLogRepository defines the interface for user activity log data access
type UserActivityLogRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, log *entity.UserActivityLog) error
	GetByID(ctx context.Context, id uint) (*entity.UserActivityLog, error)

	// User activity operations
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.UserActivityLog, error)
	LogActivity(ctx context.Context, userID uint, action, resource, ipAddress, userAgent string, metadata entity.JSONMap) error

	// List operations
	List(ctx context.Context, offset, limit int) ([]*entity.UserActivityLog, error)
	GetRecentActivity(ctx context.Context, limit int) ([]*entity.UserActivityLog, error)
}
