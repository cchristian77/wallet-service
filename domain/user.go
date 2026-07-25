package domain

import (
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FullName string
	Email    string
	Password string
}

func (User) TableName() string {
	return "users"
}
