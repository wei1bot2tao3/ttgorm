package orm

// // Predicate 谓词 查询的对象
//
//	type Predicate struct {
//		// Column 列名
//		Column string
//		// Op 运算符
//		Op string
//		// Arg 值
//		Arg any
//	}
//
// 衍生类型
type op string

const (
	opEq  op = "="
	opNot op = "NOT"
	opAnd op = "AND"
	opOR  op = "OR"
	opLT  op = "<"
	opGT  op = ">"
)

//// PredicateV1 谓词 查询的对象
//type PredicateV1 struct {
//	// Column 列名
//	c Column
//	// Op 运算符
//	op op
//	// Arg 值
//	arg any
//}

// Predicate 谓词（代表一个查询条件） 查询的对象 左边 中间 右边 做成一个复杂二叉树
type Predicate struct {
	// left 二叉树左边 查询条件左边
	left Expression
	// Op 运算符
	op op
	// Arg 值
	// 查询条件右边
	//arg any
	right Expression
}

func (p op) string() string {
	return string(p)
}

func (Predicate) expr() {}

//// EqV1 等于号
//func EqV1(column string, arg any) Predicate {
//	return Predicate{
//		left:   column,
//		op:  opEq,
//		right: arg,
//	}
//}

type Value struct {
	val any
}

func (Value) expr() {}

// And 使用C("id").Eq(12).Adn(C("name").Eq("tt))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

// Or 使用C("id").Eq(12).Or(C("name").Eq("tt))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}

// Not 是一个逻辑运算符 取反
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}
