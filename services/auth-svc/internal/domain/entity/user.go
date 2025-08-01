package entity

import "time"

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	Username        string     `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash    string     `json:"-" gorm:"not null"`
	Firstname       string     `json:"firstname"`
	Lastname        string     `json:"lastname"`
	Phone           string     `json:"phone"`
	AvatarURL       string     `json:"avatar_url"`
	IsActive        bool       `json:"is_active"`
	IsVerified      bool       `json:"is_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`

	// Relationships
	Roles        []UserRole        `json:"roles,omitempty" gorm:"foreignKey:UserID"`
	Permissions  []UserPermission  `json:"permissions,omitempty" gorm:"foreignKey:UserID"`
	ActivityLogs []UserActivityLog `json:"activity_logs,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
