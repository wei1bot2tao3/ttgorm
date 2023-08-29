package orm

import (
	"database/sql"
	"github.com/beego/beego/v2/server/web/context"
)

// Querier 用于Select语句 query 的名词形式查询
type Querier[T any] interface {
	// Get 查询单个
	Get(ctx context.Context) (*T, error)
	// GetMulti 查询 多个 multiple 多个
	GetMulti(ctx context.Context) (*[]T, error)
}

// Executor 用于INSERT，DELETE，UPDATE 执行者
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// QueryBuilder 代表SQL构造过程
type QueryBuilder interface {
	Build() (*Query, error)
}

// Query 查询参数
type Query struct {
	SQL  string
	Args []any
}

// TableName 用接口来自定义表名
type TableName interface {
	TableName() string
}
