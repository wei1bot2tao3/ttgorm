package ttOrm

import "github.com/beego/beego/v2/client/orm"

type User struct {
}

func init() {
	// 注册模型 ，数据库一张表对应到go里面一个结构体  对应到user表的结构体
	orm.RegisterModel(new(User))
	//
	orm.RegisterDriver("sqlite3", orm.DRSqlite)

	//
	orm.RegisterDataBase("default", "sqlit3", "beego.db")

}
