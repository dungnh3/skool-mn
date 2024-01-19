package models

import "github.com/dungnh3/skool-mn/internal/models/store"

type (
	Account struct {
		store.Account
	}

	StudentParent struct {
		store.StudentParent
	}

	Register struct {
		store.Register
	}

	Transaction struct {
		store.Transaction
	}
)
