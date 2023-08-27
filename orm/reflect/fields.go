package reflect

import (
	"errors"

	"reflect"
)

// InterateFields 遍历他里面的字段
func InterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持nil")
	}
	typ := reflect.TypeOf(entity)

	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	for typ.Kind() == reflect.Pointer {
		// 取他的Elem
		typ = typ.Elem()
		val = val.Elem()
	}
	for typ.Kind() != reflect.Struct {
		// 取他的Elem
		return nil, errors.New("不支持类型")
	}
	numfidld := typ.NumField()

	res := make(map[string]any, numfidld)
	for i := 0; i < numfidld; i++ {
		// 字段类型
		filedType := typ.Field(i)

		//字段的值
		fileValue := val.Field(i)
		if filedType.IsExported() {
			res[filedType.Name] = fileValue.Interface()
		} else {
			res[filedType.Name] = reflect.Zero(filedType.Type).Interface()
		}

	}
	return res, nil
}

func SetField(entity any, field string, newVal any) error {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}

	fieldVal := val.FieldByName(field)
	if fieldVal.CanSet() {
		fieldVal.Set(reflect.ValueOf(newVal))
	} else {
		return errors.New("不允许修改")
	}

	return nil
}
