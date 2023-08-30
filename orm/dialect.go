package orm

import (
	"ttgorm/orm/internal/errs"
)

var (
	DialectMySQL      Dialect = mysqlDialect{}
	DialectSQLite     Dialect = sqliteDialect{}
	DialectPostgreSQL Dialect = postgreDialect{}
)

type Dialect interface {
	// quoter 返回一个引号，引用列名，表名的引号
	quoter() byte
	//buildOnDuplicateKey构造插入冲突部分
	buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error
}

// standardSQL 其他方言继承这个接口
type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	panic("implement me")
}

// mysqlDialect 保证和 standardSQL
type mysqlDialect struct {
	standardSQL
}

func (s mysqlDialect) quoter() byte {
	return '`'
}

func (s mysqlDialect) buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	// 开始构建DUpLICAT KEY

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for k, assign := range odk.assigns {
		if k > 0 {
			b.sb.WriteByte(',')
		}
		switch a := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldsMap[a.colmun]
			if !ok {
				return errs.NewErrUnknownField(a.colmun)
			}
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.ColName)
			b.sb.WriteByte('`')
			b.sb.WriteString("=?")
			b.args = append(b.args, a.val)
		case Column:
			fd, ok := b.model.FieldsMap[a.name]
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.ColName)
			b.sb.WriteByte('`')
			b.sb.WriteString("=VALUES(")
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.ColName)
			b.sb.WriteByte('`')
			b.sb.WriteByte(')')

		default:

			return errs.NewErrUnsupportedAssignable(a.assign)
		}
	}
	return nil

}

type sqliteDialect struct {
	standardSQL
}

type postgreDialect struct {
	standardSQL
}
