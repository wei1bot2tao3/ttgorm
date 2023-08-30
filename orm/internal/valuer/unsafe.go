package valuer

import (
	"database/sql"
	"reflect"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/model"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	val   any
}

func NewUnsafeValue(model *model.Model, val any) Value {
	return unsafeValue{
		model: model,
		val:   val,
	}
}

var _ Creator = NewUnsafeValue

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	// 获取一个新的指向 T的结构体 go的数据

	// 创建一个切牌呢来存值 我先把他绑定好 因为rows.scan可以把值写进去
	var vals []any
	// 获取值的起始地址
	address := reflect.ValueOf(u.val).UnsafePointer()
	// 我得判断一下 是不是匹配的
	for _, c := range cs {

		filed, ok := u.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		// 计算偏移量 ➕起始字段的地址
		fieldAfddress := unsafe.Pointer(uintptr(address) + filed.Offset)
		val := reflect.NewAt(filed.Type, fieldAfddress)
		vals = append(vals, val.Interface())
	}
	// vals 所有的值都是 填好了
	// 现在只是翻译成元数据了 需要从元数据到 T
	err = rows.Scan(vals...)

	//  把他和tp绑定
	//使用unsafe

	return nil
}
