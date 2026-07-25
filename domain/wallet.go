package domain

import (
	"gorm.io/gorm"
)

type Wallet struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID  uint64
	Balance int64

	User *User `gorm:"foreignKey:UserID"`
}

func (Wallet) TableName() string {
	return "wallets"
}
