package domain

import "github.com/cchristian77/wallet-service/domain/enums"

type Transaction struct {
	BaseModel

	TransactionID string `gorm:"size:255;uniqueIndex"` // Idempotency-Key
	Status        enums.TransactionStatus
}

func (Transaction) TableName() string {
	return "transactions"
}
