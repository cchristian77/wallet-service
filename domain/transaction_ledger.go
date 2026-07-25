package domain

import "github.com/cchristian77/wallet-service/domain/enums"

type TransactionLedger struct {
	BaseModel

	TransactionID uint64 `gorm:"index"`
	WalletID      uint64 `gorm:"index"`
	Ledger        enums.TransactionLedgerType
	Reference     enums.TransactionLedgerReference `gorm:"index"`
	Amount        int64
}

func (TransactionLedger) TableName() string {
	return "transaction_ledgers"
}
