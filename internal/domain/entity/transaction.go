package entity

import "time"

type Transaction struct {
	ID       int64     `json:"id"`
	Sender   int64     `json:"from"`
	Receiver int64     `json:"to"`
	Amount   float64   `json:"amount"`
	Time     time.Time `json:"time"`
}
