package domain

import "github.com/cchristian77/wallet-service/domain/enums"

type TransactionLedger struct {
	BaseModel

	TransactionID uint64 `gorm:"index"` // FK → transactions.id (PK)
	WalletID      uint64
	Ledger        enums.TransactionLedgerType
	Reference     enums.TransactionLedgerReference
	Amount        int64
}

func (TransactionLedger) TableName() string {
	return "transaction_ledgers"
}
