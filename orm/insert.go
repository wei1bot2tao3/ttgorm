package orm

import (
	"context"
	"database/sql"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/model"
)

type Inserter[T any] struct {
	builder
	values []*T

	columns []string
	// 方案一
	//onConflict []Assignable
	//方案二 维持一个这个参数 buildr 模式开给他赋值
	onDuplicate *Upsert
	session     Session
}

// NewInserter 创建一个Inserter的实例 初始化好 db和builder（公共字段）
func NewInserter[T any](session Session) *Inserter[T] {
	c := session.getCore()
	return &Inserter[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		session: session,
	}
}

// Values 插入的表
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

// Columns 插入指定列
func (i *Inserter[T]) Columns(clos ...string) *Inserter[T] {
	i.columns = clos
	return i
}

// Build 构建INSTER语句
func (i Inserter[T]) Build() (*Query, error) {
	// 判断是否传入要插入的结构体

	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	i.sb.WriteString("INSERT INTO ")
	if i.model == nil {
		// 要插入结构体对应的 元数据  传入的是对应的表 在go的结构体
		m, err := i.r.Get(i.values[0])
		// model 取到的元数据 写入公共字段的 model
		i.model = m
		if err != nil {
			return nil, err
		}
	}

	//  从公共字段中 拼接表名
	i.builder.quote(i.model.TableName)

	// 一定要显示指定列的顺序，不然我我们不知道数据库中默认的顺序
	// 一定要构造 test_model
	// 开始写入数据库列名
	i.sb.WriteByte('(')
	fileds := i.model.Fields
	if len(i.columns) > 0 {
		fileds = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			// 判单是元数据中的 列名
			filed, ok := i.model.FieldsMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fileds = append(fileds, filed)
		}
	}

	// 不能遍历map 因为在go里面每异常都不一样
	for idx, filed := range fileds {
		if idx > 0 {
			i.sb.WriteByte(',')
		}

		i.quote(filed.ColName)

	}
	i.sb.WriteByte(')')
	// VALUE 开始构建 读取参数

	i.sb.WriteString(" VALUES ")

	i.args = make([]any, 0, len(i.values)*len(fileds))

	for j, v := range i.values {

		if j > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')

		for idx, field := range fileds {
			val := i.creator(i.model, v)
			if idx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			// 把参数读出来 先拿到他的反射 然后是个指针 让然后掉对应字段的值 最后转普通表达
			arg, err := val.Field(field.GoName)
			if err != nil {
				return nil, err
			}
			i.addArgs(arg)

		}
		i.sb.WriteByte(')')

	}

	// 开始构建DupLicate KEY
	if i.onDuplicate != nil {
		err := i.dialect.buildUpsert(&i.builder, i.onDuplicate)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.builder.sb.String(),
		Args: i.args,
	}, nil

}

// Assignable 标记接口
// 实现这个接口 就可以用于复制语句
type Assignable interface {
	assign()
}

type Upsert struct {
	assigns         []Assignable
	conflictColumns []string
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

// UpsertBuilder 表示开起
type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

// ConflictColumns 中间方法
func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

// Update 添加
func (o UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicate = &Upsert{
		assigns:         assigns,
		conflictColumns: o.conflictColumns,
	}
	return o.i
}

// Exec 执行 INster语句
func (i *Inserter[T]) Exec(ctx context.Context) Result {
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}

	root := i.execHandler
	for a := len(i.mdls) - 1; a >= 0; a-- {
		root = i.mdls[a](root)
	}
	res := root(ctx, &QueryContext{
		Type:    "INSERT",
		Builder: i,
		Model:   i.model,
	})
	var sqlRes sql.Result
	if res.Result != nil {
		sqlRes = res.Result.(sql.Result)
	}
	return Result{
		err: err,
		res: sqlRes,
	}
}

var _ Handler = (&Inserter[any]{}).execHandler

func (i *Inserter[T]) execHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}

	}

	res, err := i.session.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Err:    err,
		Result: res,
	}

}
