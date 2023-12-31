package v1

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"ttgorm/orm/internal/errs"
)

// Selector 定一个以查询操作的模型
type Selector[T any] struct {
	table string
	model *Model
	where []Predicate
	sb    *strings.Builder
	args  []any

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
	}
}

// Build 构建sql语句
func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = s.db.r.Registry(new(T))
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
		filename, ok := s.model.fieldsMap[exp.name]

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
func (s *Selector[T]) Form(table string) *Selector[T] {
	s.table = table
	return s
}

// Where 接收参数 中间方法 添加条件
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db

	// 在这里发起查询
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	// 确认没有数据
	if !rows.Next() {
		return nil, errs.ErrNoRows
	}

	//在这里处理结果集
	//tp:=new(T)

	// 我怎么知道你SELECT出来哪些列
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// 你就拿到了SELECT的列
	// cs 怎么处理
	// 通过cs 构造vals

	tp := new(T)

	vals := make([]any, 0, len(cs))

	valElem := make([]reflect.Value, 0, len(cs))

	for _, c := range cs {

		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}
		// 是*int
		val := reflect.New(fd.typ)
		//vals = append(vals, val.Elem())

		vals = append(vals, val.Interface())
		//记得调用
		valElem = append(valElem, val.Elem())
	}

	//第一个问题：类型要匹配
	//第二个问题：顺序要匹配

	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}
	tpValueElem := reflect.ValueOf(tp).Elem()
	for i, c := range cs {
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)

		}
		tpValueElem.FieldByName(fd.GOName).Set(valElem[i])

	}
	// x想办法把vals塞进去 tp 里面

	return tp, err
}

func (s *Selector[T]) GetMulti(ctx context.Context) (*[]T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err

	}
	db := s.db.db
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if rows.Next() {

	}
	return nil, err
}
