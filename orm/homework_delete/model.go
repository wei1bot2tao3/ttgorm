package homework_delete

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
	NumField := typ.NumField()
	FieldMap := make(map[string]*field, NumField)
	for i := 0; i < NumField; i++ {
		filedType := typ.Field(i)
		FieldMap[filedType.Name] = &field{
			colName: underscoreName(filedType.Name),
		}
	}
	return &model{
		tableName: typ.Name(),
		fields:    FieldMap,
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
