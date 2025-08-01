package repository

import (
	"context"
	"errors"

	"github.com/taititans/bitzap/auth-svc/internal/domain/entity"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"gorm.io/gorm"
)

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername gets user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

// List gets users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// Count counts total users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error
	return count, err
}

// Search searches users
func (r *userRepository) Search(ctx context.Context, query string, offset, limit int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Where("email ILIKE ? OR username ILIKE ? OR firstname ILIKE ? OR lastname ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// GetWithRoles gets user with roles
func (r *userRepository) GetWithRoles(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetWithPermissions gets user with permissions
func (r *userRepository) GetWithPermissions(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Permissions").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetWithActivityLogs gets user with activity logs
func (r *userRepository) GetWithActivityLogs(ctx context.Context, id uint, limit int) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("ActivityLogs", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(limit)
		}).
		First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin updates user's last login time
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

// VerifyEmail verifies user's email
func (r *userRepository) VerifyEmail(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_verified":       true,
			"email_verified_at": gorm.Expr("NOW()"),
		}).Error
}

// VerifyPhone verifies user's phone
func (r *userRepository) VerifyPhone(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).
		Update("phone_verified_at", gorm.Expr("NOW()")).Error
}
