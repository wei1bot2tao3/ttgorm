package orm

// Aggregate 聚合函数 常见聚合函数AVG("age") ,SUM("age")
type Aggregate struct {
	// fn 聚合函数名
	fn string
	// arg 字段名 goName
	arg   string
	alias string
}

func (a Aggregate) selectable() {

}

// Avg AVG函数  平均值
func Avg(col string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: col,
	}
}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}

// Sum sum 函数
func Sum(col string) Aggregate {
	return Aggregate{
		fn:  "SUM",
		arg: col,
	}
}

// Max MAX函数
func Max(col string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		arg: col,
	}
}

// Count Count函数
func Count(col string) Aggregate {
	return Aggregate{
		fn:  "COUNT",
		arg: col,
	}
}

// Min Mint函数
func Min(col string) Aggregate {
	return Aggregate{
		fn:  "MIN",
		arg: col,
	}
}
