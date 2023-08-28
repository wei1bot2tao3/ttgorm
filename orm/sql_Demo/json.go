package sql_Demo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JsonColumn[T any] struct {
	Val T

	Valid bool
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	//
	var bs []byte
	switch data := src.(type) {
	case string:
		bs = []byte(data)
	case []byte:
		bs = data
	case nil:
		// 说明数据里没数据
		return nil
	default:
		return errors.New("不支持的类型")

	}
	err := json.Unmarshal(bs, &j.Val)
	fmt.Println(err)
	if err == nil {
		j.Valid = true
	}
	return err
}
