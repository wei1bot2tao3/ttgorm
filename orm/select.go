package orm

import (
	"context"
	"fmt"
	"ttgorm/orm/internal/errs"
)

// Selectable
type Selectable interface {
	selectable()
}

// Selector 定一个以查询操作的模型
type Selector[T any] struct {
	builder
	table TableReference
	//model   *model.Model
	where []Predicate
	//sb      *strings.Builder
	//args    []any
	columns []Selectable
	//db      *DB
	session Session
}

func NewSelector[T any](session Session) *Selector[T] {
	c := session.getCore()
	return &Selector[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		session: session,
	}
}

// Build 构建sql语句
func (s *Selector[T]) Build() (*Query, error) {
	if s.model == nil {
		var err error
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}

	s.sb.WriteString("SELECT ")
	if err := s.BuilderColumns(); err != nil {
		return nil, err
	}

	s.sb.WriteString(" FROM ")
	if err := s.buildTable(s.table); err != nil {
		return nil, err
	}

	//if s.table == "" {
	//
	//	//我怎么把表名字拿到
	//	s.sb.WriteByte('`')
	//	s.sb.WriteString(s.model.TableName)
	//	s.sb.WriteByte('`')
	//} else {
	//	//segs := strings.Split(s.table, ".")
	//	//for i, v := range segs {
	//	//	s.sb.WriteByte('`')
	//	//	s.sb.WriteString(v)
	//	//	s.sb.WriteByte('`')
	//	//	if i < len(segs)-1 {
	//	//		s.sb.WriteByte('.')
	//	//	}
	//	//}
	//	// 加不加引号？
	//	s.sb.WriteByte('`')
	//	s.sb.WriteString(s.model.TableName)
	//	s.sb.WriteByte('`')
	//}

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
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

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil

}

func (s *Selector[T]) buildTable(table TableReference) error {
	switch t := table.(type) {
	case nil:
		// 这是代表完全没有调用 FROM，也就是最普通的形态
		s.quote(s.model.TableName)
	case Table:
		m, err := s.r.Get(t.entity)
		if err != nil {
			return err
		}
		s.quote(m.TableName)
		if t.alias != "" {
			s.sb.WriteString(" AS ")
			s.quote(t.alias)
		}
	case Join:
		s.sb.WriteByte('(')
		// 构造左边
		err := s.buildTable(t.left)
		if err != nil {
			return err
		}
		s.sb.WriteByte(' ')
		// 构造中间
		s.sb.WriteString(t.typ)
		s.sb.WriteByte(' ')
		// 构造右边
		err = s.buildTable(t.right)
		if err != nil {
			return err
		}

		if len(t.using) > 0 {
			s.sb.WriteString(" USING (")
			// 拼接 USING (xx, xx)
			for i, col := range t.using {
				if i > 0 {
					s.sb.WriteByte(',')
				}
				err = s.buildColumn(Column{name: col})
				if err != nil {
					return err
				}
			}
			s.sb.WriteByte(')')
		}

		if len(t.on) > 0 {
			s.sb.WriteString(" ON ")
			p := t.on[0]
			for i := 1; i < len(t.on); i++ {
				p = p.And(t.on[i])
			}
			if err = s.buildExpression(p); err != nil {
				return err
			}
		}

		s.sb.WriteByte(')')
	default:
		return errs.NewErrUnsupportedTable(table)
	}
	return nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		// 没有指定列
		s.sb.WriteByte('*')
		return nil
	}

	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch c := col.(type) {
		case Column:
			err := s.buildColumn(c)
			if err != nil {
				return err
			}
		case Aggregate:
			// 聚合函数名
			s.sb.WriteString(c.fn)
			s.sb.WriteByte('(')
			err := s.buildColumn(Column{name: c.arg})
			if err != nil {
				return err
			}
			s.sb.WriteByte(')')
			// 聚合函数本身的别名
			if c.alias != "" {
				s.sb.WriteString(" AS `")
				s.sb.WriteString(c.alias)
				s.sb.WriteByte('`')
			}
		case RawExpr:
			s.sb.WriteString(c.raw)
			s.addArg(c.args...)
		}
	}

	return nil
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
		exp.alias = ""
		err := s.buildColumn(exp)
		if err != nil {
			return err
		}

	//右边
	case Value:
		// 把参数放进去
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	case RawExpr:
		s.sb.WriteByte('(')
		s.sb.WriteString(exp.raw)
		s.addArg(exp.args...)
		s.sb.WriteByte(')')

	default:
		return fmt.Errorf("orm : 不支持的分支")
	}

	return nil

}

// BuilderColumns 构造列 聚合函数
func (s *Selector[T]) BuilderColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
	}
	if len(s.columns) > 0 {
		for i, v := range s.columns {
			if i > 0 {
				s.sb.WriteByte(',')
			}

			switch c := v.(type) {
			case Column:
				err := s.buildColumn(c)
				if err != nil {
					return err
				}
			case Aggregate:
				s.sb.WriteString(c.fn)
				s.sb.WriteString(`(`)
				err := s.buildColumn(Column{
					name: c.arg,
				})
				if err != nil {
					return err
				}
				s.sb.WriteString(`)`)
				if c.alias != "" {
					s.sb.WriteString(" AS `")
					s.sb.WriteString(c.alias)
					s.sb.WriteByte('`')
				}

			case RawExpr:
				s.sb.WriteString(c.raw)
				s.addArg(c.args...)
			}

		}
	}

	return nil
}

func (s *Selector[T]) addArg(vals ...any) {
	if len(vals) == 0 {
		return
	}
	if s.args == nil {
		s.args = make([]any, 0, 16)
	}
	s.args = append(s.args, vals...)

}

// From 中间方法 要原封不懂返回 Selector 这个是添加表名
func (s *Selector[T]) From(table TableReference) *Selector[T] {
	s.table = table
	return s
}

// SelectV1 直接拼接
//func (s *Selector[T]) SelectV1(columns ...string) *Selector[T] {
//	s.columns = columns
//	return s
//}

// Select
func (s *Selector[T]) Select(columns ...Selectable) *Selector[T] {
	s.columns = columns
	return s
}

// Where 接收参数 中间方法 添加条件
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.session, s.core, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   s.model,
	})

	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (s *Selector[T]) GetMulti(ctx context.Context) (*[]T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err

	}

	rows, err := s.session.queryContext(ctx, q.SQL, q.Args...)
	if rows.Next() {

	}
	return nil, err
}
