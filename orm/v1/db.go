package v1

import (
	"database/sql"
)

// DB sql.DB的装饰器
type DB struct {
	r  *registry
	db *sql.DB
}

type DBOption func(db *DB)

// Open 注册一个实例 返回一个 同时注册一个元数据实例 放在db中返回
func Open(driver string, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)

}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:  &registry{},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
