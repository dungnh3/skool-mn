package mysql

import (
	"context"
	"github.com/dungnh3/skool-mn/config"
	"github.com/dungnh3/skool-mn/internal/models"
	"github.com/dungnh3/skool-mn/internal/repositories"
	l "github.com/dungnh3/skool-mn/pkg/log"
	"gorm.io/gorm"
)

type Queries struct {
	db     *gorm.DB
	cfg    *config.Config
	logger l.Logger
}

func (q *Queries) UpdateRegister(ctx context.Context, r *models.Register) error {
	//TODO implement me
	panic("implement me")
}

func New(db *gorm.DB, cfg *config.Config) *Queries {
	q := Queries{
		db:     db,
		cfg:    cfg,
		logger: l.New().Named("database"),
	}
	return &q
}

func (q *Queries) Ping() error {
	sqlDB, err := q.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (q *Queries) WithTx(tx *gorm.DB) *Queries {
	nr := *q
	nr.db = tx
	return &nr
}

func (q *Queries) Transaction(txFunc func(repositories.Repository) error) (err error) {
	// start new transaction
	tx := q.db.Begin()
	defer func() {
		p := recover()
		switch {
		case p != nil:
			execErr := tx.Rollback().Error
			if execErr != nil {
				q.logger.Error("error exec rollback", l.Error(execErr))
			}
			panic(p) // re-throw panic after Rollback
		case err != nil:
			execErr := tx.Rollback().Error // err is non-nil; don't change it
			if execErr != nil {
				q.logger.Error("error exec rollback", l.Error(execErr))
			}
		default:
			err = tx.Commit().Error // err is nil; if Commit returns error update err
		}
	}()
	return txFunc(q.WithTx(tx))
}
