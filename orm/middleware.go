package orm

import "ttgorm/orm/model"

// QueryContext 代表查询上下文
type QueryContext struct {
	// 查询类型，标记增删改查
	Type string

	// 代表的是查询本身
	Builder QueryBuilder

	Model *model.Model
}
