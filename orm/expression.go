package orm

// Expression 代表语句  标记接口
type Expression interface {
	expr()
}

// RowExpr 代表了原始表达式
// 意味着 ORM对他不会有任何处理
type RowExpr struct {
	raw  string
	args []any
}

func (r RowExpr) selectable() {

}

// Raw 创建一个RawExpr
func Raw(expr string, args ...any) RowExpr {
	return RowExpr{
		raw:  expr,
		args: args,
	}
}
func (r RowExpr) expr() {

}

// AsPredicate  让他本身作为一个Predicate
func (r RowExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
