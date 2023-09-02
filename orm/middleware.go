package orm

import (
	"context"
	"ttgorm/orm/model"
)

// QueryContext 代表查询上下文
type QueryContext struct {
	// 查询类型，标记增删改查
	Type string

	// 代表的是查询本身
	Builder QueryBuilder

	Model *model.Model
}

// QueryResult   代表查询结果
type QueryResult struct {

	// Result 在不同的查询类型下是不同的
	// SELECT可以是 *T
	Result any
	Err    error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
