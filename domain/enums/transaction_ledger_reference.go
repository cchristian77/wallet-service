package enums

type TransactionLedgerReference string

const (
	TransactionLedgerRefDisbursement TransactionLedgerReference = "DISBURSEMENT"
	TransactionLedgerRefTopUp        TransactionLedgerReference = "TOP_UP"
	TransactionLedgerRefTransfer     TransactionLedgerReference = "TRANSFER"
	TransactionLedgerRefWithdrawal   TransactionLedgerReference = "WITHDRAWAL"
)

func (tlr TransactionLedgerReference) String() string {
	return string(tlr)
}
