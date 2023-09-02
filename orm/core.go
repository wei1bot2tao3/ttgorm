package orm

import (
	"context"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/internal/valuer"
	"ttgorm/orm/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator

	r    model.Registry
	mdls []Middleware
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

func getHandler[T any](ctx context.Context, session Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	rows, err := session.queryContext(ctx, q.SQL, q.Args...)

	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next() {
		return &QueryResult{
			Err: errs.ErrNoRows,
		}
	}
	//

	// 我现在拿到了查询的数据 我要把 我的数据库的数据 通过元数据 翻译成go的数据

	tp := new(T)
	// 怎么把 	valuer.Value() 和tp 关联在一起 使用一个工厂模式

	val := c.creator(c.model, tp)
	err = val.SetColumns(rows)

	// 接口定义好改造上层，用
	return &QueryResult{
		Err:    err,
		Result: tp,
	}

}

func execHandler(ctx context.Context, session Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
			Result: Result{
				err: err,
			},
		}

	}

	res, err := session.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Err: err,
		Result: Result{
			err: err,
			res: res,
		},
	}

}

func exec(ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {

	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return execHandler(ctx, sess, c, qc)
	}
	for i := len(c.mdls) - 1; i >= 0; i-- {
		root = c.mdls[i](root)
	}

	return root(ctx, qc)

}
