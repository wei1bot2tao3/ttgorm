package orm

import "context"

// RawQuerier 实现原生查询
type RawQuerier[T any] struct {
	core
	session Session
	sql     string
	aegs    []any
}

func (s *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  s.sql,
		Args: s.aegs,
	}, nil
}

func RawQuery[T any](query string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  query,
		aegs: args,
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
func get[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {

	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, sess, c, qc)
	}
	for i := len(c.mdls) - 1; i >= 0; i-- {
		root = c.mdls[i](root)
	}

	return root(ctx, qc)

}

//
//func (s *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
//	var err error
//	s.model, err = s.r.Get(new(T))
//	if err != nil {
//		return nil, err
//	}
//
//	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
//		return getHandler[T](ctx,s.session,s.core,qc)
//	}
//	for i := len(s.mdls) - 1; i >= 0; i-- {
//		root = s.mdls[i](root)
//	}
//
//	res := root(ctx, &QueryContext{
//		Type:    "RAW",
//		Builder: s,
//		Model:   s.model,
//	})
//
//	if res.Result != nil {
//		return res.Result.(*T), res.Err
//	}
//	return nil, res.Err
//
//}