package repositories

import (
	"context"
	"github.com/dungnh3/skool-mn/internal/models"
)

type Repository interface {
	Ping() error
	Transaction(txFunc func(Repository) error) error
	StudentParentRepository
	RegisterRepository
	TransactionRepository
}

type StudentParentRepository interface {
	ListStudents(ctx context.Context, parentId string) ([]*models.Account, error)
}

type RegisterRepository interface {
	CreateRegister(ctx context.Context, r *models.Register) error
	ConfirmFromTeacher(ctx context.Context, registerId string) error
	RejectFromTeacher(ctx context.Context, registerId string) error
	ConfirmFromParent(ctx context.Context, registerId string) error
	WaitingFromParent(ctx context.Context, registerId string) error
	CancelFromParent(ctx context.Context, registerId string) error
	GetConfirmedRegisterLatest(ctx context.Context, studentId string) (*models.Register, error)
	GetLeftRegisterLatest(ctx context.Context, studentId string) (*models.Register, error)
	StudentLeaveClass(ctx context.Context, registerId string) error
	StudentOutSchool(ctx context.Context, registerId string) error
	LockRegister(ctx context.Context, id string) (*models.Register, error)
	UpdateRegister(ctx context.Context, r *models.Register) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
}
