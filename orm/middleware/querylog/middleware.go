package querylog

import (
	"context"
	"log"
	"ttgorm/orm"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("sql: %s, args: %v", query, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				// 要考虑记录下来吗？
				log.Println("构造SQL出错", err)
				return &orm.QueryResult{
					Err: err,
				}

			}
			m.logFunc(q.SQL, q.Args)
			res := next(ctx, qc)
			return res

		}

	}
}
