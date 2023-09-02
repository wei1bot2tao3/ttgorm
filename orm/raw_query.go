package orm

import (
	"context"
	"database/sql"
)

// RawQuerier 实现原生查询
type RawQuerier[T any] struct {
	core
	session Session
	sql     string
	aegs    []any
}

func (s RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  s.sql,
		Args: s.aegs,
	}, nil
}

func RawQuery[T any](sess Session, query string, args ...any) *RawQuerier[T] {
	c := sess.getCore()
	return &RawQuerier[T]{
		sql:     query,
		aegs:    args,
		session: sess,
		core:    c,
	}
}

func (s *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.r.Get(new(T))
	if err != nil {
		return nil, err
	}

	res := get[T](ctx, s.session, s.core, &QueryContext{
		Type:    "RAW",
		Builder: s,
		Model:   s.model,
	})

	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err

}

//func (i RawQuerier[T]) Exec(ctx context.Context) Result {
//	var err error
//	i.model, err = i.r.Get(new(T))
//	if err != nil {
//		return Result{
//			err: err,
//		}
//	}
//
//	res := exec(ctx, i.session, i.core, &QueryContext{
//		Type: "RAW",
//		Builder: i,
//		Model: i.model,
//	} )
//	var sqlRes sql.Result
//	if res.Result != nil {
//		sqlRes = res.Result.(sql.Result)
//	}
//	return Result{
//		err: err,
//		res: sqlRes,
//	}
//}

//

func (i RawQuerier[T]) Exec(ctx context.Context) Result {
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}

	res := exec(ctx, i.session, i.core, &QueryContext{
		Type:    "RAW",
		Builder: i,
		Model:   i.model,
	})
	// var t *T
	// if val, ok := res.Result.(*T); ok {
	// 	t = val
	// }
	// return t, res.Err
	var sqlRes sql.Result
	if res.Result != nil {
		sqlRes = res.Result.(sql.Result)
	}
	return Result{
		err: res.Err,
		res: sqlRes,
	}
}
