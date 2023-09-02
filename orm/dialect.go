package orm

import (
	"ttgorm/orm/internal/errs"
)

var (
	DialectMySQL      Dialect = mysqlDialect{}
	DialectSQLite     Dialect = sqlite3Dialect{}
	DialectPostgreSQL Dialect = postgreDialect{}
)

// Dialect 定义了与数据库方言相关的方法。
type Dialect interface {
	// quoter 返回一个引号，引用列名，表名的引号
	quoter() byte
	//buildOnDuplicateKey构造插入冲突部分
	buildUpsert(b *builder, upsert *Upsert) error
}

// standardSQL 其他方言继承这个接口
type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildUpsert(b *builder, upsert *Upsert) error {
	panic("implement me")
}

// mysqlDialect 保证和 standardSQL
type mysqlDialect struct {
	standardSQL
}

func (s mysqlDialect) quoter() byte {
	return '`'
}

func (s mysqlDialect) buildUpsert(b *builder, upsert *Upsert) error {
	// 开始构建DUpLICAT KEY

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for k, assign := range upsert.assigns {
		if k > 0 {
			b.sb.WriteByte(',')
		}
		switch a := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldsMap[a.colmun]
			if !ok {
				return errs.NewErrUnknownField(a.colmun)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.args = append(b.args, a.val)
		case Column:
			fd, ok := b.model.FieldsMap[a.name]
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sb.WriteByte(')')

		default:

			return errs.NewErrUnsupportedAssignable(a.assign)
		}
	}
	return nil

}

type sqlite3Dialect struct {
	standardSQL
}

func (s sqlite3Dialect) quoter() byte {
	return '`'
}
func (s sqlite3Dialect) buildUpsert(b *builder, upsert *Upsert) error {
	// 开始构建DUpLICAT KEY
	b.sb.WriteString("ON CONFLICT(")
	for i, col := range upsert.conflictColumns {
		if i > 0 {
			b.sb.WriteByte(',')

		}
		err := b.buildColumn(Column{
			name: col,
		})
		if err != nil {
			return err
		}
	}
	b.sb.WriteString(")")
	//
	b.sb.WriteString(" DO UPDATE SET ")
	for idx, assign := range upsert.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch a := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldsMap[a.colmun]
			// 字段不对，或者说列不对
			if !ok {
				return errs.NewErrUnknownField(a.colmun)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.addArgs(a.val)

		case Column:
			fd, ok := b.model.FieldsMap[a.name]
			// 字段不对，或者说列不对
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=excluded.")
			b.quote(fd.ColName)
		default:
			return errs.NewErrUnsupportedAssignable(assign)
		}
	}

	return nil
}

type postgreDialect struct {
	standardSQL
}
