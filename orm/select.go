package orm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"ttgorm/orm/internal/errs"
	"unsafe"
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
	// 在这里发起查询 查询的结果放到 rows里面
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	// 确认没有数据 判断rows里也没有数据
	if !rows.Next() {
		return nil, errs.ErrNoRows
	}

	//在这里处理结果集    处理rows
	// 我怎么知道你SELECT出来哪些列
	//rows.Columns()：用于获取结果集的列名列表，返回一个字符串切片，其中包含查询结果的列名。
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// 你就拿到了SELECT的列
	// cs 怎么处理
	// 通过cs 构造vals
	//new(T)可以创建一个指向T类型零值的指针
	tp := new(T)
	// 创建一个切片来存放val 查询出来的值
	vals := make([]any, 0, len(cs))
	//存放了查询结果的反射值
	valElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		//  s.model.columnMap 存放是数据库的列名 判断查询出来的是否符合结构体
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}
		// 存在返回一个 field 一个字段结构体 里面有列名 字段类型
		// 是*int
		// 通过字段的映射创建一个新的反射值 ，本质上是 查询到的数据 写到这个里 所有创建一个新的反射值
		val := reflect.New(fd.typ)
		//vals = append(vals, val.Elem())
		// 把val的指针 存进去 后面靠rows.Scan 来把值填进去
		vals = append(vals, val.Interface())
		//记得调用 存的是映射值 可以靠这个获取指针和值
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

func (s *Selector[T]) GetV1(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	db := s.db.db

	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)

	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errs.ErrNoRows
	}

	// 我现在拿到了查询的数据 我要把 我的数据库的数据 通过元数据 翻译成go的数据
	// 拿到了 查询结果的列名
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// 获取一个新的指向 T的结构体 go的数据

	tp := new(T)
	// 创建一个切牌呢来存值 我先把他绑定好 因为rows.scan可以把值写进去
	vals := make([]any, 0, len(cs))
	valElem := make([]reflect.Value, 0, len(cs))
	// 我得判断一下 是不是匹配的
	for _, c := range cs {

		filed, ok := s.model.columnMap[c]
		if !ok {
			return nil, errors.New("这个列和数据库有一个不匹配")
		}
		//reflect.New 返回的是reflect.Value类型的值 根据他的typ创建一个指针 是tpvalue的
		val := reflect.New(filed.typ)
		// 如果匹配了 转换次Interface 因为 Sacn 只收any
		vals = append(vals, val.Interface())
		valElem = append(valElem, val.Elem())
	}
	// vals 所有的值都是 填好了
	// 现在只是翻译成元数据了 需要从元数据到 T
	err = rows.Scan(vals...)
	if err != nil {
		return nil, err
	}
	//  把他和tp绑定
	tpValueElem := reflect.ValueOf(tp).Elem()
	for i, c := range cs {
		filed, ok := s.model.columnMap[c]
		if !ok {
			return nil, errors.New("这个列和数据库有一个不匹配")
		}
		// 现在是把元数据的值给 T 首先获取t的值ELem 然后根据Byname+goName 进行复制 set必须传的是  reflect.Value
		tpValueElem.FieldByName(filed.GOName).Set(valElem[i])
	}

	return tp, nil
}

func (s *Selector[T]) GetV2(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	db := s.db.db
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errs.ErrNoRows
	}

	// 我现在拿到了查询的数据 我要把 我的数据库的数据 通过元数据 翻译成go的数据
	// 拿到了 查询结果的列名
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// 获取一个新的指向 T的结构体 go的数据

	tp := new(T)
	// 创建一个切牌呢来存值 我先把他绑定好 因为rows.scan可以把值写进去
	var vals []any
	// 获取值的起始地址
	address := reflect.ValueOf(tp).UnsafePointer()
	// 我得判断一下 是不是匹配的
	for _, c := range cs {

		filed, ok := s.model.columnMap[c]
		if !ok {
			return nil, errors.New("这个列和数据库有一个不匹配")
		}

		// 计算偏移量 ➕起始字段的地址
		fieldAfddress := unsafe.Pointer(uintptr(address) + filed.Offset)
		val := reflect.NewAt(filed.typ, fieldAfddress)
		vals = append(vals, val.Interface())
	}
	// vals 所有的值都是 填好了
	// 现在只是翻译成元数据了 需要从元数据到 T
	err = rows.Scan(vals...)

	//  把他和tp绑定
	//使用unsafe

	return tp, nil
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
