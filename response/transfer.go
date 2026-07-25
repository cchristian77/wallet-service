package response

type Transfer struct {
	TransactionID string `json:"transaction_id"`
	FromBalance   int64  `json:"from_balance"`
	ToBalance     int64  `json:"to_balance"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
}
