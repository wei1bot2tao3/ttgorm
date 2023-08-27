package v1

import (
	"reflect"
	"ttgorm/orm/internal/errs"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	// 列名
	colName string
}

func parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}

	typ = typ.Elem()
	//for typ.Kind() == reflect.Pointer {
	//	typ = typ.Elem()
	//}

	numFiled := typ.NumField()
	fieldMap := make(map[string]*field, numFiled)
	for i := 0; i < numFiled; i++ {
		filedType := typ.Field(i)
		fieldMap[filedType.Name] = &field{
			colName: underscoreName(filedType.Name),
		}
	}

	return &model{

		tableName: underscoreName(typ.Name()),

		fields: fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
