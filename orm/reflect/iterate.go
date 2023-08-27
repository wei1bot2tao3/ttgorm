package reflect

import "reflect"

func IterateArry(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	// 告诉你多长
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Cap(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil

}
func IterateSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	// 告诉你多长
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Cap(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil

}

// IterateMap keys ,values ,err
func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	// 告诉你多长
	resValues := make([]any, 0, val.Len())
	resKeys := make([]any, 0, val.Len())
	Keys := val.MapKeys()
	for _, key := range Keys {
		value := val.MapIndex(key)
		resKeys = append(resKeys, key.Interface())
		resValues = append(resValues, value.Interface())

	}
	return resKeys, resValues, nil

}

func IterateMapV2(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	// 告诉你多长
	resValues := make([]any, 0, val.Len())
	resKeys := make([]any, 0, val.Len())
	itr := val.MapRange()
	for itr.Next() {
		resKeys = append(resKeys, itr.Key().Interface())
		resValues = append(resValues, itr.Value().Interface())
	}

	return resKeys, resValues, nil

}