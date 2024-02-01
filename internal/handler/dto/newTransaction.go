package dto

type NewTransactionDTO struct {
	To     int64   `json:"to"`
	Amount float64 `json:"amount"`
}
