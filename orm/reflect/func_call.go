package reflect

import (
	"fmt"
	"reflect"
)

func IterateFunc(entity any) map[string]FuncInfo {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		// 遍历方法
		method := typ.Method(i)
		// 获取方法的func
		fn := method.Func
		// 获取方法输入变量的数量
		numIn := fn.Type().NumIn()

		// 传入变量的类型
		input := make([]reflect.Type, 0, numIn)
		// 传入变量的值
		inputValues := make([]reflect.Value, 0, numIn)
		//
		inputValues = append(inputValues, reflect.ValueOf(entity))
		input = append(input, reflect.TypeOf(entity))

		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			input = append(input, fnInType)
			inputValues = append(inputValues, reflect.Zero(fnInType))
		}

		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))

		}
		resValues := fn.Call(inputValues)
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}
		fmt.Println("out是啥", output)
		res[method.Name] = FuncInfo{
			Name:       method.Name,
			InputTypes: input,
			OutTypes:   output,
			Result:     result,
		}
	}
	return res
}

type FuncInfo struct {
	Name       string
	InputTypes []reflect.Type
	OutTypes   []reflect.Type
	Result     []any
}
