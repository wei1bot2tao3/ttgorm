package orm

import (
	"context"
	"database/sql"
	"ttgorm/orm/internal/errs"
	"ttgorm/orm/internal/valuer"
	"ttgorm/orm/model"
)

var (
	_ Session = &Tx{}
	_ Session = &DB{}
)

// DB sql.DB的装饰器
type DB struct {
	db *sql.DB

	core
}

type DBOption func(db *DB)

// Open 注册一个实例 返回一个 同时注册一个元数据实例 放在db中返回
// driver 表示数据库驱动的名称  dataSourceNam表示数据库连接的数据源名称。它是一个字符串，包含了连接数据库所需的信息，如主机名、端口号、数据库名称、用户名和密码
// opts ...DBOption：表示可选的数据库选项
func Open(driver string, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)

}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			r:       model.NewRegistry(),
			creator: valuer.NewUnsafeValue,
			dialect: DialectMySQL,
		},
		db: db,
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

func DbWithMiddlewares(mdls ...Middleware) DBOption {
	return func(db *DB) {
		db.mdls = mdls
	}
}

// DBWithDialect 允许使用方言
func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBUserReflect() DBOption {
	return func(db *DB) {
		db.creator = valuer.NewReflectValue
	}
}

// BeginTx 它用于在数据库连接上开始一个事务 opts 是一个指向 sql.TxOptions 类型的指针，用于设置事务的选项
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)

	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
	}, nil
}

type txKey struct {
}

// BeginTXV2  没有事务我给你开 扩散 可以引入一个done的标记为
func (db *DB) BeginTXV2(ctx context.Context, otps *sql.TxOptions) (context.Context, *Tx, error) {
	val := ctx.Value(txKey{})
	tx, ok := val.(*Tx)
	// 存在一个事务，并且这个事务没有被提交或者回滚
	if ok && !tx.done {
		return ctx, tx, nil
	}

	tx, err := db.BeginTx(ctx, otps)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, txKey{}, tx)
	return ctx, tx, nil

}

// 要求前面的人一定要开好事务
// func (db *DB) BeginTxV3(ctx context.Context,
// 	opts *sql.TxOptions) (*Tx, error){
// 	val := ctx.Value(txKey{})
// 	tx, ok := val.(*Tx)
// 	if ok {
// 		return tx, nil
// 	}
// 	return nil, errors.New("没有开事务")
// }

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

func (db *DB) getCore() core {
	return db.core

}

// DoTx 实现闭包API
func (db *DB) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, otps *sql.TxOptions) (err error) {
	tx, err := db.BeginTx(ctx, otps)
	if err != nil {
		return nil
	}
	panicked := true

	defer func() {
		if panicked || err != nil {
			// 回滚
			e := tx.Rollback()
			err = errs.NewErrFailedToRollbackTx(err, e, panicked)
		} else {
			// 提交
			err = tx.Commit()

		}

	}()

	err = fn(ctx, tx)
	panicked = false
	return err

}
