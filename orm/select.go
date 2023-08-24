package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	sb    *strings.Builder
	args  []any
}

// Build 构建sql语句
func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	sb := s.sb
	sb.WriteString("SELECT * FORM ")
	if s.table == "" {

		//我怎么把表名字拿到
		var t T
		sb.WriteByte('`')
		sb.WriteString(reflect.TypeOf(t).Name())
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
		sb.WriteString(s.table)
		sb.WriteByte('`')
	}

	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])

		}
		//在这里处理p
		s.buildExpression(p)

	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil

}

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
		s.sb.WriteString(exp.name)
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
