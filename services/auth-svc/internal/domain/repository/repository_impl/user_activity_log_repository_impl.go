package repository

import (
	"context"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"gorm.io/gorm"
)

// userActivityLogRepository implements UserActivityLogRepository
type userActivityLogRepository struct {
	db *gorm.DB
}

// NewUserActivityLogRepository creates a new user activity log repository
func NewUserActivityLogRepository(db *gorm.DB) repository.UserActivityLogRepository {
	return &userActivityLogRepository{db: db}
}

// Create creates a new user activity log
func (r *userActivityLogRepository) Create(ctx context.Context, log *entity.UserActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByID gets user activity log by ID
func (r *userActivityLogRepository) GetByID(ctx context.Context, id uint) (*entity.UserActivityLog, error) {
	var log entity.UserActivityLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

// GetByUserID gets user activity logs by user ID
func (r *userActivityLogRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.UserActivityLog, error) {
	var logs []*entity.UserActivityLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

// LogActivity logs user activity
func (r *userActivityLogRepository) LogActivity(ctx context.Context, userID uint, action, resource, ipAddress, userAgent string, metadata entity.JSONMap) error {
	log := &entity.UserActivityLog{
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Metadata:  metadata,
	}
	return r.db.WithContext(ctx).Create(log).Error
}

// List gets user activity logs with pagination
func (r *userActivityLogRepository) List(ctx context.Context, offset, limit int) ([]*entity.UserActivityLog, error) {
	var logs []*entity.UserActivityLog
	err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

// GetRecentActivity gets recent activity logs
func (r *userActivityLogRepository) GetRecentActivity(ctx context.Context, limit int) ([]*entity.UserActivityLog, error) {
	var logs []*entity.UserActivityLog
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
