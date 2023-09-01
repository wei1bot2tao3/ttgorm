package valuer

import (
	"database/sql"
	"ttgorm/orm/model"

	"reflect"
	"ttgorm/orm/internal/errs"
)

type reflectValue struct {
	model *model.Model
	//val   any
	val reflect.Value
}

func NewReflectValue(model *model.Model, val any) Value {
	return reflectValue{
		model: model,
		val:   reflect.ValueOf(val).Elem(),
	}
}

func (r reflectValue) Field(name string) (any, error) {
	// 检查name是否合法
	//_,ok:=r.val.Type().FieldByName(name)
	//if !ok{
	//	return nil,errors.New("错了")
	//}
	return r.val.FieldByName(name).Interface(), nil
}

var _ Creator = NewReflectValue

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	//在这里处理结果集    处理rows
	// 我怎么知道你SELECT出来哪些列
	//rows.Columns()：用于获取结果集的列名列表，返回一个字符串切片，其中包含查询结果的列名。
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	// 你就拿到了SELECT的列
	// cs 怎么处理
	// 通过cs 构造vals
	//new(T)可以创建一个指向T类型零值的指针

	//tp := new(r.val)
	// 创建一个切片来存放val 查询出来的值
	vals := make([]any, 0, len(cs))
	//存放了查询结果的反射值
	valElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		//  s.model.columnMap 存放是数据库的列名 判断查询出来的是否符合结构体
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 存在返回一个 field 一个字段结构体 里面有列名 字段类型
		// 是*int
		// 通过字段的映射创建一个新的反射值 ，本质上是 查询到的数据 写到这个里 所有创建一个新的反射值
		val := reflect.New(fd.Type)
		//vals = append(vals, val.Elem())
		// 把val的指针 存进去 后面靠rows.Scan 来把值填进去
		vals = append(vals, val.Interface())
		//记得调用 存的是映射值 可以靠这个获取指针和值
		valElem = append(valElem, val.Elem())
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	// 处理元数据到go 结构体
	tpValueElem := reflect.ValueOf(r.val).Elem()

	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)

		}
		tpValueElem.FieldByName(fd.GoName).Set(valElem[i])

	}
	// x想办法把vals塞进去 tp 里面

	return err
}
