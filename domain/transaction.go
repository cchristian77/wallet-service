package domain

import (
	"github.com/cchristian77/wallet-service/domain/enums"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TransactionID string `gorm:"size:255;uniqueIndex"` // Idempotency-Key
	Status        enums.TransactionStatus
}

func (Transaction) TableName() string {
	return "transactions"
}
