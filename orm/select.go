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
	table string
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
	if s.table == "" {

		//我怎么把表名字拿到
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		//segs := strings.Split(s.table, ".")
		//for i, v := range segs {
		//	s.sb.WriteByte('`')
		//	s.sb.WriteString(v)
		//	s.sb.WriteByte('`')
		//	if i < len(segs)-1 {
		//		s.sb.WriteByte('.')
		//	}
		//}
		// 加不加引号？
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	}

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
		err := s.BuilderColumn(exp)
		if err != nil {
			return err
		}

	//右边
	case Value:
		// 把参数放进去
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	case RowExpr:
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
				err := s.BuilderColumn(c)
				if err != nil {
					return err
				}
			case Aggregate:
				s.sb.WriteString(c.fn)
				s.sb.WriteString(`(`)
				err := s.BuilderColumn(Column{
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

			case RowExpr:
				s.sb.WriteString(c.raw)
				s.addArg(c.args...)
			}

		}
	}

	return nil
}

func (s *Selector[T]) BuilderColumn(c Column) error {

	filename, ok := s.model.FieldsMap[c.name]

	if !ok {
		// 传入错误 或者列不队
		return errs.NewErrUnknownField(c.name)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(filename.ColName)
	s.sb.WriteByte('`')
	if c.alias != "" {
		s.sb.WriteString(" AS `")
		s.sb.WriteString(c.alias)
		s.sb.WriteByte('`')
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
func (s *Selector[T]) From(table string) *Selector[T] {
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

func getHandler[T any](ctx context.Context, session Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	rows, err := session.queryContext(ctx, q.SQL, q.Args...)

	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next() {
		return &QueryResult{
			Err: errs.ErrNoRows,
		}
	}
	//

	// 我现在拿到了查询的数据 我要把 我的数据库的数据 通过元数据 翻译成go的数据

	tp := new(T)
	// 怎么把 	valuer.Value() 和tp 关联在一起 使用一个工厂模式

	val := c.creator(c.model, tp)
	err = val.SetColumns(rows)

	// 接口定义好改造上层，用
	return &QueryResult{
		Err:    err,
		Result: tp,
	}

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
