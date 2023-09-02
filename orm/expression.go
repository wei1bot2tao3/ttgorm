package orm

// Expression 代表语句  标记接口
type Expression interface {
	expr()
}

// RawExpr 代表了原始表达式
// 意味着 ORM对他不会有任何处理
type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectable() {

}

// Raw 创建一个RawExpr
func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}
func (r RawExpr) expr() {

}

// AsPredicate  让他本身作为一个Predicate
func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
