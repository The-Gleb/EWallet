package maketransaction_usecase

type MakeTransactionDTO struct {
	From   int64
	To     int64
	Amount float64
}
