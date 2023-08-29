package valuer

import (
	"database/sql"
	"reflect"
	"ttgorm/orm/internal/errs"
	v1 "ttgorm/orm/model"
	"unsafe"
)

type unSafevalue struct {
	model *v1.Model
	// 对应t的指针
	val any
}

func NewunSafevalue(model *v1.Model, val any) Valuer {
	return reflectValue{
		model: model,
		val:   val,
	}
}

var _ Creator = NewunSafevalue

func (r unSafevalue) SetColumns(rows *sql.Rows) error {

	//在这里处理结果集
	//tp:=new(T)
	// 我怎么知道你SELECT出来哪些列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	var vals []any

	//起始地址
	address := reflect.ValueOf(r.val).UnsafePointer()
	for _, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 起始地址加偏移量
		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}

	err = rows.Scan(vals...)
	return err
}
