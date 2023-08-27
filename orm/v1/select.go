package v1

import (
	"context"
	"fmt"
	"strings"
	"ttgorm/orm/internal/errs"
)

type Selector[T any] struct {
	table string
	model *model
	where []Predicate
	sb    *strings.Builder
	args  []any
}

// Build 构建sql语句
func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = parseModel(new(T))
	if err != nil {
		return nil, err
	}
	sb := s.sb
	sb.WriteString("SELECT * FORM ")
	if s.table == "" {

		//我怎么把表名字拿到
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')
	} else {
		//segs := strings.Split(s.table, ".")
		//for i, v := range segs {
		//	sb.WriteByte('`')
		//	sb.WriteString(v)
		//	sb.WriteByte('`')
		//	if i < len(segs)-1 {
		//		sb.WriteByte('.')
		//	}
		//}
		// 加不加引号？
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')
	}

	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])

		}
		//在这里处理p
		err := s.buildExpression(p)
		if err != nil {
			return nil, err
		}

	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil

}

// buildExpression 构建SQL的条件
func (s *Selector[T]) buildExpression(expr Expression) error {

	switch exp := expr.(type) {
	case nil:
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		// 中间
		s.sb.WriteString(exp.op.string())
		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

	// 左边
	case Column:
		s.sb.WriteByte('`')
		filename, ok := s.model.fields[exp.name]

		if !ok {
			// 传入错误 或者列不队
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteString(filename.colName)
		s.sb.WriteByte('`')

	//右边
	case Value:
		// 把参数放进去
		s.sb.WriteByte('?')
		s.addArg(exp.val)

	default:
		return fmt.Errorf("orm : 不支持的分支")
	}

	return nil

}

func (s *Selector[T]) addArg(val any) *Selector[T] {
	if s.args == nil {
		s.args = make([]any, 0, 16)
	}
	s.args = append(s.args, val)
	return s
}

// Form 中间方法 要原封不懂返回 Selector 这个是添加表名
func (s *Selector[T]) Form(tabel string) *Selector[T] {
	s.table = tabel
	return s
}

// Where 接收参数 中间方法 添加条件
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) (*[]T, error) {
	//TODO implement me
	panic("implement me")
}
