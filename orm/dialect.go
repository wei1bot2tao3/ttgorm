package orm

import (
	"strings"
	"ttgorm/orm/internal/errs"
)

type Dialect interface {
	// quoter 返回一个引号，引用列名，表名的引号
	quoter() byte
	//buildOnDuplicateKey构造插入冲突部分
	buildOnDuplicateKey(sb strings.Builder, odk *OnDuplicateKey) error
}

// standardSQL 其他方言继承这个接口
type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildOnDuplicateKey(sb strings.Builder, odk *OnDuplicateKey) error {

}

//  mysqlDialect 保证和 standardSQL
type mysqlDialect struct {
	standardSQL
}

func (s mysqlDialect) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s mysqlDialect) buildOnDuplicateKey(sb strings.Builder, odk *OnDuplicateKey) error {
	// 开始构建DUpLICAT KEY

	sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for k, assign := range odk.assigns {
		if k > 0 {
			sb.WriteByte(',')
		}
		switch a := assign.(type) {
		case Assignment:
			fd, ok := m.FieldsMap[a.colmun]
			if !ok {
				return errs.NewErrUnknownField(a.colmun)
			}
			sb.WriteByte('`')
			sb.WriteString(fd.ColName)
			sb.WriteByte('`')
			sb.WriteString("=?")
			args = append(args, a.val)
		case Column:
			fd, ok := m.FieldsMap[a.name]
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			sb.WriteByte('`')
			sb.WriteString(fd.ColName)
			sb.WriteByte('`')
			sb.WriteString("=VALUES(")
			sb.WriteByte('`')
			sb.WriteString(fd.ColName)
			sb.WriteByte('`')
			sb.WriteByte(')')

		default:

			return errs.NewErrUnsupportedAssignable(a.assign)
		}
	}
	return nil

}

type sqliteDialect struct {
	standardSQL
}
