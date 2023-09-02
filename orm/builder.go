package orm

import (
	"strings"
	"ttgorm/orm/internal/errs"
)

// builder 公共的抽象 需要什么放进去什么
type builder struct {
	sb   strings.Builder
	args []any
	core
	quoter byte
}

// quote name
func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

// buildColumn 构建列 把列名➕到sb里

func (b *builder) buildColumn(c Column) error {
	switch table := c.table.(type) {
	case nil:
		fd, ok := b.model.FieldsMap[c.name]
		// 字段不对，或者说列不对
		if !ok {
			return errs.NewErrUnknownField(c.name)
		}
		b.quote(fd.ColName)
		if c.alias != "" {
			b.sb.WriteString(" AS ")
			b.quote(c.alias)
		}
	case Table:
		m, err := b.r.Get(table.entity)
		if err != nil {
			return err
		}
		fd, ok := m.FieldsMap[c.name]
		if !ok {
			return errs.NewErrUnknownField(c.name)
		}
		if table.alias != "" {
			b.quote(table.alias)
			b.sb.WriteByte('.')
		}
		b.quote(fd.ColName)
		if c.alias != "" {
			b.sb.WriteString(" AS ")
			b.quote(c.alias)
		}
	default:
		return errs.NewErrUnsupportedTable(table)
	}
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
