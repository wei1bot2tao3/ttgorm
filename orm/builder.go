package orm

import (
	"strings"
	"ttgorm/orm/model"
)

// builder 公共的抽象 需要什么放进去什么
type builder struct {
	sb      strings.Builder
	args    []any
	model   model.Model
	dialect Dialect
	quoter  byte
}

// quote name
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString("name")
	b.sb.WriteByte(b.quoter)
}

func (b *builder)  {

}
