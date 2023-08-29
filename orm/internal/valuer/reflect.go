package valuer

import (
	"database/sql"
	"reflect"
	"ttgorm/orm/internal/errs"
	v1 "ttgorm/orm/model"
)

type reflectValue struct {
	model *v1.Model
	// 对应t的指针
	val any
}

func NewReflectValue(model *v1.Model, val any) Valuer {
	return reflectValue{
		model: model,
		val:   val,
	}
}

var _ Creator = NewReflectValue

func (r reflectValue) SetColumns(rows *sql.Rows) error {

	//在这里处理结果集
	//tp:=new(T)
	// 我怎么知道你SELECT出来哪些列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	// 你就拿到了SELECT的列
	// cs 怎么处理
	// 通过cs 构造vals

	vals := make([]any, 0, len(cs))
	valElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {

		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 是*int
		val := reflect.New(fd.Typ)
		//vals = append(vals, val.Elem())

		vals = append(vals, val.Interface())
		//记得调用
		valElem = append(valElem, val.Elem())
	}

	//第一个问题：类型要匹配
	//第二个问题：顺序要匹配

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}
	tpValueElem := reflect.ValueOf(r.val).Elem()
	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)

		}
		tpValueElem.FieldByName(fd.GOName).Set(valElem[i])

	}
	// x想办法把vals塞进去 tp 里面

	return err
}
