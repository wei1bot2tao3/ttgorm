package valuer

import (
	"database/sql"
	v1 "ttgorm/orm/model"
)

type Valuer interface {
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(model *v1.Model, entity any) Valuer
