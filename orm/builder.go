package orm

import (
	"strings"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/model"
)

// builder 公共的抽象 需要什么放进去什么
type builder struct {
	sb      strings.Builder
	args    []any
	model   *model.Model
	dialect Dialect
	quoter  byte
}

// quote name
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

// buildColumn 构建列 把列名➕到sb里
func (b *builder) buildColumn(name string) error {
	fd, ok := b.model.FieldsMap[name]
	if !ok {
		return errs.NewErrUnknownField(name)
	}
	b.quote(fd.ColName)
	return nil

}

func (b *builder) addArgs(args ...any) error {
	if b.args == nil {
		// 很少有函数超过八个参数
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
	return nil
}
