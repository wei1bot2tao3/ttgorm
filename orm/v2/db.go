package v1

type DB struct {
	r *registry
}

type DBOption func(db *DB)

// NewDB 注册一个实例 返回一个 同时注册一个元数据实例 放在db中返回
func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: newRegistry(),
	}
	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func MustNewDB(opts ...DBOption) *DB {
	res, err := NewDB(opts...)
	if err != nil {
		panic(err)
	}
	return res
}
