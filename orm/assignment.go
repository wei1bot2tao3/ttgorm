package orm

// Assignment "assignment" 表示将一个值赋给一个变量或列的操作
type Assignment struct {
	colmun string
	val    any
}

func Assign(col string, val any) Assignment {
	return Assignment{
		colmun: col,
		val:    val,
	}
}

func (Assignment) assign() {

}
