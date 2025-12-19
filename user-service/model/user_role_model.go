package model

import "time"

type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	RoleID    uint      `json:"roles_id" gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Role      Role      `gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
