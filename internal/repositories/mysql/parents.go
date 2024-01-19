package mysql

import (
	"context"
	"database/sql"
	"github.com/dungnh3/skool-mn/internal/models"
)

const listStudentsQuery = `
SELECT a.*
FROM accounts a
WHERE a.id IN (SELECT sp.student_id FROM student_parents sp WHERE sp.parent_id = @parent_id)
`

func (q *Queries) ListStudents(
	ctx context.Context, parentId string,
) ([]*models.Account, error) {
	var accounts []*models.Account
	tx := q.db.WithContext(ctx).Debug().Raw(
		listStudentsQuery,
		sql.Named("parent_id", parentId),
	).Find(&accounts)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
