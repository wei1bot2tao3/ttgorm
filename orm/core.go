package orm

import (
	"ttgorm/orm/internal/valuer"
	"ttgorm/orm/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator

	r model.Registry
}
