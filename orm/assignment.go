package orm

type Assignment struct {
	colmun string
	val    any
}

func Assgin(col string, val any) Assignment {
	return Assignment{
		colmun: col,
		val:    val,
	}
}

func (Assignment) assign() {

}
