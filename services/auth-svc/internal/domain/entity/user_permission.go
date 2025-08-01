package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserPermission struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Resource  string         `json:"resourcce" gorm:"not null"`
	Action    string         `json:"action" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (UserPermission) TableName() string {
	return "user_permissions"
}
