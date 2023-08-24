package orm

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
}

func (s *Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	sb.WriteString("SELECT * FORM ")
	if s.table == "" {

		//我怎么把表名字拿到
		var t T
		sb.WriteByte('`')
		sb.WriteString(reflect.TypeOf(t).Name())
		sb.WriteByte('`')
	} else {
		segs := strings.Split(s.table, ".")
		for i, v := range segs {
			sb.WriteByte('`')
			sb.WriteString(v)
			sb.WriteByte('`')
			if i < len(segs)-1 {
				sb.WriteByte('.')
			}
		}
		// 加不加引号？
		//sb.WriteString(s.table)
	}
	sb.WriteByte(';')
	return &Query{
		SQL: sb.String(),
	}, nil

}

// Form 中间方法 要原封不懂返回 Selector
func (s *Selector[T]) Form(tabel string) *Selector[T] {
	s.table = tabel
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) (*[]T, error) {
	//TODO implement me
	panic("implement me")
}
