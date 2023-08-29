package v1

import (
	"strings"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/model"
)

type Deleter[T any] struct {
	// 存储 sql语句
	sb    strings.Builder
	m     *model.Model
	table string
	where []Predicate
	args  []any
	r     *model.registry
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.m, _ = d.r.Registry(new(T))
	d.sb.WriteString("DELETE FROM ")

	if d.table == "" {

		table := model.underscoreName(d.m.TableName)
		d.sb.WriteByte('`')
		d.sb.WriteString(table)
		d.sb.WriteByte('`')
	} else {

		d.sb.WriteString(d.table)

	}
	// 处理 Where

	if len(d.where) > 0 {
		pred := d.where[0]
		d.sb.WriteString(" WHERE ")
		for i := 1; i < len(d.where)-1; i++ {
			pred = pred.And(d.where[i])
		}

		err := d.buildDeleteExpression(pred)
		if err != nil {
			return nil, err
		}

	}

	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

// BuildDeletepExpression 构建语句
func (d *Deleter[T]) buildDeleteExpression(expression Expression) error {

	switch exp := expression.(type) {
	case nil:
	case Predicate:
		if exp.left != nil {

			_, ok := exp.left.(Predicate)
			if ok {
				d.sb.WriteByte('(')
			}
			err := d.buildDeleteExpression(exp.left)
			if err != nil {
				return err
			}
			if ok {
				d.sb.WriteByte(')')
			}
		}
		// 中间
		d.sb.WriteString(exp.op.string())

		// 右边

		if exp.right != nil {
			_, ok := exp.right.(Predicate)
			if ok {
				d.sb.WriteByte('(')
			}
			err := d.buildDeleteExpression(exp.right)
			if err != nil {
				return err
			}
			if ok {
				d.sb.WriteByte(')')
			}
		}

	case Column:

		filename, ok := d.m.FieldsMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		d.sb.WriteByte('`')
		columnname := model.underscoreName(filename.ColName)
		d.sb.WriteString(columnname)
		d.sb.WriteByte('`')
	case Value:
		d.sb.WriteByte('?')
		d.addArg(exp.val)
	}

	return nil
}
func (d *Deleter[T]) addArg(value any) *Deleter[T] {
	d.args = make([]any, 0, 16)
	d.args = append(d.args, value)
	return d
}

// From accepts Model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = model.underscoreName(table)
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}

func NewDeleter[T any]() *Deleter[T] {
	return &Deleter[T]{}
}
