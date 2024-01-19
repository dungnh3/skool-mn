package models

const (
	AccountsTable     = "accounts"
	TransactionsTable = "transactions"
	RegistersTable    = "registers"
)

func (a *Account) TableName() string {
	return AccountsTable
}

func (a *Transaction) TableName() string {
	return TransactionsTable
}

func (a *Register) TableName() string {
	return RegistersTable
}
