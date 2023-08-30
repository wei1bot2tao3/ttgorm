package orm

// Column  列名
type Column struct {
	name  string
	alias string
}

// C 组装列名
func C(name string) Column {
	return Column{
		name: name,
	}
}
func (c Column) assign() {

}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}
func (c Column) expr() {}

// Eq 更方便= C("id").Eq(arg) 链式调用
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: c.valueOf(arg),
	}
}

func (c Column) valueOf(arg any) Expression {
	switch val := arg.(type) {
	case Expression:
		return val
	default:
		return Value{
			val: val,
		}
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opLT,
		right: Value{
			val: arg,
		},
	}
}

func (c Column) selectable() {

}
