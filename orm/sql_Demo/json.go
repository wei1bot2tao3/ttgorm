package sql_Demo

import (
	"database/sql/driver"
	"fmt"
)

type JsonColumn[T any] struct {
	Val T

	Valid bool
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (j *JsonColumn[T]) Scan(state fmt.ScanState, verb rune) error {
	//TODO implement me
	panic("implement me")
}
