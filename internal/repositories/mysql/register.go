package mysql

import (
	"context"
	"github.com/dungnh3/skool-mn/internal/models"
	"github.com/dungnh3/skool-mn/internal/models/store"
	"gorm.io/gorm/clause"
	"time"
)

func (q *Queries) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	return q.db.WithContext(ctx).Debug().Model(&models.Transaction{}).Create(tx).Error
}

func (q *Queries) CreateRegister(ctx context.Context, r *models.Register) error {
	return q.db.WithContext(ctx).Debug().Model(&models.Register{}).Create(r).Error
}

func (q *Queries) LockRegister(ctx context.Context, id string) (*models.Register, error) {
	var r models.Register
	tx := q.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&models.Register{}).Where("id = ?", id)
	return &r, tx.First(&r).Error
}

func (q *Queries) ConfirmFromTeacher(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.Confirmed, "updated_at": time.Now()}).Error
}

func (q *Queries) RejectFromTeacher(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.Rejected, "updated_at": time.Now()}).Error
}

func (q *Queries) ConfirmFromParent(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.Done, "updated_at": time.Now()}).Error
}

func (q *Queries) WaitingFromParent(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.Waiting, "updated_at": time.Now()}).Error
}

func (q *Queries) CancelFromParent(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.Cancelled, "updated_at": time.Now()}).Error
}

func (q *Queries) StudentLeaveClass(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.StudentLeftClass, "updated_at": time.Now()}).Error
}

func (q *Queries) StudentOutSchool(ctx context.Context, registerId string) error {
	return q.db.WithContext(ctx).Model(&models.Register{}).
		Where("id = ?", registerId).
		Updates(map[string]interface{}{"status": store.StudentOutSchool, "updated_at": time.Now()}).Error
}

func (q *Queries) GetConfirmedRegisterLatest(ctx context.Context, studentId string) (*models.Register, error) {
	var r models.Register
	return &r, q.db.WithContext(ctx).Model(&models.Register{}).
		Where("student_id = ? AND status = ?", studentId, store.Confirmed).
		Order("created_at desc").Take(&r).Error
}

func (q *Queries) GetLeftRegisterLatest(ctx context.Context, studentId string) (*models.Register, error) {
	var r models.Register
	return &r, q.db.WithContext(ctx).Model(&models.Register{}).Debug().
		Where("student_id = ? AND status = ?", studentId, store.StudentLeftClass).
		Order("created_at desc").Take(&r).Error
}
