package enums

type TransactionLedgerType string

const (
	TransactionLedgerDebit  TransactionLedgerType = "DEBIT"
	TransactionLedgerCredit TransactionLedgerType = "CREDIT"
)

func (tlt TransactionLedgerType) String() string {
	return string(tlt)
}
