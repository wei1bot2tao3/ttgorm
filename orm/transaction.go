package orm

import "database/sql"

type Tx struct {
	tx *sql.Tx
}

// Commit 提交一个事务
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollbac 回滚一个事务
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
