package valuer

import (
	"database/sql"
	"ttgorm/orm/model"
)

type Value interface {
	// SetColumns 设置新的值 把元数据和返回值结合
	SetColumns(rows *sql.Rows) error
}

type Creator func(model *model.Model, entity any) Value
