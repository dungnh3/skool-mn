package repositories

import (
	"github.com/dungnh3/skool-mn/internal/models/store"
)

type StudentParentFilter struct {
	ParentId  string
	ArrStatus []store.ObjectStatus
}
