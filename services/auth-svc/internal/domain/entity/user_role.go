package entity

import "time"

type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
