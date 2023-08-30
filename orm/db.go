package orm

import (
	"database/sql"
	"ttgorm/orm/internal/valuer"
	"ttgorm/orm/model"
)

// DB sql.DB的装饰器
type DB struct {
	r       model.Registry
	db      *sql.DB
	creator valuer.Creator
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
		r:       model.NewRegistry(),
		db:      db,
		creator: valuer.NewUnsafeValue,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBUserReflect() DBOption {
	return func(db *DB) {
		db.creator = valuer.NewReflectValue
	}
}
