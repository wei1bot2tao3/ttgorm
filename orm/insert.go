package orm

import (
	"reflect"
	"strings"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/model"
)

type Inserter[T any] struct {
	values  []*T
	db      *DB
	columns []string
	// 方案一
	//onConflict []Assignable
	//方案二 维持一个这个参数 buildr 模式开给他赋值
	onDuplicate *OnDuplicateKey
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

// Values 插入
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

// Cloumns 插入指定列
func (i *Inserter[T]) Cloumns(clos ...string) *Inserter[T] {
	i.columns = clos
	return i
}

// Build 构建Insert语句
func (i Inserter[T]) BuildBackup() (*Query, error) {

	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	// 拿元数据
	m, err := i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}

	// 拼接表名
	sb.WriteByte('`')
	sb.WriteString(m.TableName)
	sb.WriteByte('`')
	// 一定要显示指定列的顺序，不然我我们不知道数据库中默认的顺序
	// 一定要构造 test_model
	sb.WriteByte('(')
	fileds := m.Fields
	if len(i.columns) > 0 {
		fileds = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			filed, ok := m.FieldsMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fileds = append(fileds, filed)
		}
	}

	// 不能遍历map 因为在go里面每异常都不一样
	//
	for idx, filed := range fileds {
		if idx > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('`')
		sb.WriteString(filed.ColName)
		sb.WriteByte('`')

	}
	sb.WriteByte(')')
	// VALUE 开始构建 读取参数

	sb.WriteString(" VALUES ")

	args := make([]any, 0, len(i.values)*len(fileds))

	for j, val := range i.values {

		if j > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('(')
		for idx, field := range fileds {

			if idx > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte('?')
			// 把参数读出来 先拿到他的反射 然后是个指针 让然后掉对应字段的值 最后转普通表达
			arg := reflect.ValueOf(val).Elem().FieldByName(field.GoName).Interface()
			args = append(args, arg)

		}
		sb.WriteByte(')')

	}

	// 开始构建DUpLICAT KEY
	if i.onDuplicate != nil {
		sb.WriteString(" ON DUPLICATE KEY UPDATE ")
		for k, assign := range i.onDuplicate.assigns {
			if k > 0 {
				sb.WriteByte(',')
			}
			switch a := assign.(type) {
			case Assignment:
				fd, ok := m.FieldsMap[a.colmun]
				if !ok {
					return nil, errs.NewErrUnknownField(a.colmun)
				}
				sb.WriteByte('`')
				sb.WriteString(fd.ColName)
				sb.WriteByte('`')
				sb.WriteString("=?")
				args = append(args, a.val)
			case Column:
				fd, ok := m.FieldsMap[a.name]
				if !ok {
					return nil, errs.NewErrUnknownField(a.name)
				}
				sb.WriteByte('`')
				sb.WriteString(fd.ColName)
				sb.WriteByte('`')
				sb.WriteString("=VALUES(")
				sb.WriteByte('`')
				sb.WriteString(fd.ColName)
				sb.WriteByte('`')
				sb.WriteByte(')')

			default:

				return nil, errs.NewErrUnsupportedAssignable(a.assign)
			}
		}
	}

	sb.WriteByte(';')

	return &Query{
		SQL:  sb.String(),
		Args: args,
	}, nil

}

func (i Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	// 拿元数据
	m, err := i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}

	// 拼接表名
	sb.WriteByte('`')
	sb.WriteString(m.TableName)
	sb.WriteByte('`')
	// 一定要显示指定列的顺序，不然我我们不知道数据库中默认的顺序
	// 一定要构造 test_model
	sb.WriteByte('(')
	fileds := m.Fields
	if len(i.columns) > 0 {
		fileds = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			filed, ok := m.FieldsMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fileds = append(fileds, filed)
		}
	}

	// 不能遍历map 因为在go里面每异常都不一样
	for idx, filed := range fileds {
		if idx > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('`')
		sb.WriteString(filed.ColName)
		sb.WriteByte('`')

	}
	sb.WriteByte(')')
	// VALUE 开始构建 读取参数

	sb.WriteString(" VALUES ")

	args := make([]any, 0, len(i.values)*len(fileds))

	for j, val := range i.values {

		if j > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('(')
		for idx, field := range fileds {

			if idx > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte('?')
			// 把参数读出来 先拿到他的反射 然后是个指针 让然后掉对应字段的值 最后转普通表达
			arg := reflect.ValueOf(val).Elem().FieldByName(field.GoName).Interface()
			args = append(args, arg)

		}
		sb.WriteByte(')')

	}

	// 开始构建DUpLICAT KEY
	if i.onDuplicate != nil {

	}

	sb.WriteByte(';')

	return &Query{
		SQL:  sb.String(),
		Args: args,
	}, nil

}

// Assignable 标记接口
// 实现这个接口 就可以用于复制语句
type Assignable interface {
	assign()
}

//func (i *Inserter[T]) OnDuplicateKey(assigns ...Assignable) *Inserter[T] {
//	i.onDuplicate = assigns
//}

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

// OnDuplicateKeyBuilder 表示开起
type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}
type OnDuplicateKey struct {
	assigns []Assignable
}

// Update
func (o OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicate = &OnDuplicateKey{
		assigns: assigns,
	}
	return o.i
}
