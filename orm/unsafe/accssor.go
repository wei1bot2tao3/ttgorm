package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnSafeAccessor struct {
	fileds  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnSafeAccessor(entity any) *UnSafeAccessor {
	typ := reflect.TypeOf(entity).Elem()
	nuFiled := typ.NumField()
	fileds := make(map[string]FieldMeta, nuFiled)
	for i := 0; i < nuFiled; i++ {
		fd := typ.Field(i)
		fileds[fd.Name] = FieldMeta{
			// 存放偏移量
			Offset: fd.Offset,
			// 存放类型
			typ: fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnSafeAccessor{
		fileds: fileds,
		// 计算起始地址   uintptr
		//address:val.UnsafeAddr(), 为啥不用他 他不是稳定的 在垃圾回收前 后会变
		//使用Pointer
		address: val.UnsafePointer(),
	}
}

func (a *UnSafeAccessor) Filed(filed string) (any, error) {
	// 起始地址+ 字段偏移量
	fd, ok := a.fileds[filed]
	if !ok {
		return nil, errors.New("非法字段")
	}

	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 如果不知道类型就这么读
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
	// 如果知道类型就这么读
	//return *(*int)(fdAddress), nil
}

func (a *UnSafeAccessor) Set(filed string, val any) error {
	// 起始地址+ 字段偏移量
	fd, ok := a.fileds[filed]
	if !ok {
		return errors.New("非法字段")
	}
	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	// 假设你知道确切类型
	//*(*int)(fdAddress) = val.(int)

	//不知道确切类型
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

// FieldMeta 存放偏移量
type FieldMeta struct {
	// uintptr 是一个数字 代表内存地址
	Offset uintptr
	typ    reflect.Type
}
