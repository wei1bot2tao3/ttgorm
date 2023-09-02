package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"ttgorm/orm/internal/errs"
)

//func TestSelector_Build(t *testing.T) {
//	db := memoryDB(t)
//
//	testCasses := []struct {
//		name string
//
//		builder    QueryBuilder
//		wantQuerry *Query
//		wantErr    error
//		wantRes    *TestModel
//	}{
//		{
//			name:    "select no FROM",
//			builder: NewSelector[TestModel](db),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model`;",
//				Args: nil,
//			},
//		},
//		{
//			name:    "FROM",
//			builder: NewSelector[TestModel](db).From("test_model"),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model`;",
//				Args: nil,
//			},
//		},
//		{
//			name:    "empty FROM ",
//			builder: NewSelector[TestModel](db).From(""),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model`;",
//				Args: nil,
//			},
//		},
//		//{
//		//	name:    "FROM db ",
//		//	builder: (&Selector[TestModel]{}).FROM("test_db.test_model"),
//		//	wantQuerry: &Query{
//		//		SQL:  "SELECT * FROM `test_db`.`test_model`;",
//		//		Args: nil,
//		//	},
//		//},
//		{
//			name:    "where",
//			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18)),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model` WHERE `id`=?;",
//				Args: []any{18},
//			},
//		},
//		{
//			name:    "where",
//			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18).And(C("Id").Eq(11))),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model` WHERE (`id`=?)AND(`id`=?);",
//				Args: []any{18, 11},
//			},
//		},
//		{
//			name:    "not",
//			builder: NewSelector[TestModel](db).Where(Not(C("Id").Eq(18))),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model` WHERE NOT(`id`=?);",
//				Args: []any{18},
//			},
//		},
//		{
//			name:    "where",
//			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18).Or(C("Id").Eq(11))),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model` WHERE (`id`=?)OR(`id`=?);",
//				Args: []any{18, 11},
//			},
//		},
//		{
//			name:    "where",
//			builder: NewSelector[TestModel](db).Where(),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model`;",
//				Args: nil,
//			},
//		},
//		{
//			name:    "where",
//			builder: NewSelector[TestModel](db).Where(Not(C("jkd").Eq(18))),
//			wantQuerry: &Query{
//				SQL:  "SELECT * FROM `test_model`;",
//				Args: nil,
//			},
//			wantErr: errs.NewErrUnknownField("jkd"),
//		},
//	}
//
//	for _, tc := range testCasses {
//		t.Run(tc.name, func(t *testing.T) {
//			q, err := tc.builder.Build()
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			assert.Equal(t, tc.wantQuerry, q)
//		})
//	}
//}

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

// memoryDB  SQLite3 数据库驱动创建一个内存数据库，并返回对该数据库的引用
func memoryDB(t *testing.T, opts ...DBOption) *DB {
	db, err := Open("sqlite3",
		"file:test.db?cache=shared&mode=memory",
		// 仅仅用于单元测试，不会发起真的查询
		opts...)
	require.NoError(t, err)
	return db
}

func TestSelector_Get(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDb)
	require.NoError(t, err)
	fmt.Println(mock)
	// 对应与
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))
	// 对应与 弄rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "Jerry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("abc", "Tom", "18", "Jerry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
	testCasses := []struct {
		name string
		s    *Selector[TestModel]

		wantErr error
		wantRes *TestModel
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no roes",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errs.ErrNoRows,
		},

		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		//{
		//	name: "scan error",
		//	s:    NewSelector[TestModel](db).Where(C("Id").LT(1)),
		//	wantRes: &TestModel{
		//		Id:        1,
		//		FirstName: "Tom",
		//		Age:       18,
		//		LastName:  &sql.NullString{Valid: true, String: "Jerry"},
		//	},
		//	wantErr: errs.ErrNoRows,
		//},
	}

	for _, tc := range testCasses {
		t.Run(tc.name, func(t *testing.T) {

			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

type ResultTEste struct {
	ID   int
	Name string
}

func TestEml(t *testing.T) {
	// 模拟查询结果
	vals := []interface{}{1, "John"}

	// 创建结果对象
	tp := ResultTEste{}

	// 获取结果对象的反射值
	tpValue := reflect.ValueOf(&tp).Elem()

	// 字段名
	fd := struct{ GOName string }{GOName: "N"}

	// 在结果对象中查找指定字段的反射值
	fieldValue := tpValue.FieldByName(fd.GOName)

	// 检查字段是否存在
	if fieldValue.IsValid() {
		// 检查字段的类型是否匹配
		if fieldValue.Type().AssignableTo(reflect.TypeOf(vals[0])) {
			// 设置字段的值
			fieldValue.Set(reflect.ValueOf(vals[0]))
		} else {
			fmt.Println("字段类型不匹配")
		}
	} else {
		fmt.Println("字段不存在")
	}

	fmt.Println("ID:", tp.Name) // 输出结果：ID: 1
}

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCass := []struct {
		name    string
		q       QueryBuilder
		wantErr error
		wantRes *Query
	}{
		{
			name: "multiple columns ",
			q:    NewSelector[TestModel](db).Select(C("FirstName"), C("LastName")),
			wantRes: &Query{
				SQL: "SELECT `first_name`,`last_name` FROM `test_model`;",
			},
		},

		{
			name: "multiple columns ",
			q:    NewSelector[TestModel](db),
			wantRes: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},

		{
			name: "multiple columns ",
			q:    NewSelector[TestModel](db).Select(C("test_model")),
			wantRes: &Query{
				SQL: "SELECT `first_name`,`last_name` FROM `test_model`;",
			},
			wantErr: errs.NewErrUnknownField("test_model"),
		},

		{
			name: "AVG ",
			q:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantRes: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "SUM ",
			q:    NewSelector[TestModel](db).Select(Sum("Age")),
			wantRes: &Query{
				SQL: "SELECT SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "max ",
			q:    NewSelector[TestModel](db).Select(Max("Age")),
			wantRes: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "Min",
			q:    NewSelector[TestModel](db).Select(Min("Age")),
			wantRes: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},
		{
			name: "Min  error ",
			q:    NewSelector[TestModel](db).Select(Min("Invalid")),
			wantRes: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
			wantErr: errs.NewErrUnknownField("Invalid"),
		},

		{
			name: "Count",
			q:    NewSelector[TestModel](db).Select(Count("Age")),
			wantRes: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},

		{
			name: "Min AND MAX",
			q:    NewSelector[TestModel](db).Select(Min("Age"), Sum("Age")),
			wantRes: &Query{
				SQL: "SELECT MIN(`age`),SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "Raw",
			q:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantRes: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},

		{
			name: "Raw expression",
			q:    NewSelector[TestModel](db).Where(Raw("`age` < ?", 18).AsPredicate()),
			wantRes: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` < ?);",
				Args: []any{18},
			},
		},

		{
			name: "column alias",
			q:    NewSelector[TestModel](db).Select(C("FirstName").As("my_name"), C("LastName")),
			wantRes: &Query{
				SQL: "SELECT `first_name` AS `my_name`,`last_name` FROM `test_model`;",
			},
		},

		{
			name: "avg alias",
			q:    NewSelector[TestModel](db).Select(Avg("FirstName").As("avg_name"), C("LastName")),
			wantRes: &Query{
				SQL: "SELECT AVG(`first_name`) AS `avg_name`,`last_name` FROM `test_model`;",
			},
		},

		{
			name: "avg alias",
			q:    NewSelector[TestModel](db).Where(C("Id").As("myid").Eq(1)),
			wantRes: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id`=?;",
				Args: []any{1},
			},
		},

		{
			name: "Raw expression",
			q:    NewSelector[TestModel](db).Where(C("Id").Eq(Raw("`age`+?", 1))),
			wantRes: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id`=(`age`+?);",
				Args: []any{1},
			},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				fmt.Println(err)
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSelector_Join(t *testing.T) {
	db := memoryDB(t)
	type Order struct {
		Id        int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId  int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		s         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "specify table",
			s:    NewSelector[Order](db).From(TableOf(&OrderDetail{})),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order_detail`;",
			},
		},
		{
			name: "join-suing",
			s: func() QueryBuilder {
				t1 := TableOf(&Order{})
				t2 := TableOf(&OrderDetail{})
				t3 := t1.Join(t2).Using("UsingCol1", "UsingCol2")
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` JOIN `order_detail` USING (`using_col1`,`using_col2`));",
			},
		},
		{
			name: "right join",
			s: func() QueryBuilder {
				t1 := TableOf(&Order{})
				t2 := TableOf(&OrderDetail{})
				t3 := t1.RightJoin(t2).Using("UsingCol1", "UsingCol2")
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` RIGHT JOIN `order_detail` USING (`using_col1`,`using_col2`));",
			},
		},
		{
			name: "join-on",
			s: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.Join(t2).On(t1.C("Id").Eq(t2.C("OrderId")))
				return NewSelector[Order](db).From(t3)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id`=`t2`.`order_id`);",
			},
		},
		{
			name: "join table",
			s: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := t1.Join(t2).On(t1.C("Id").Eq(t2.C("OrderId")))
				t4 := TableOf(&Item{}).As("t4")
				t5 := t3.Join(t4).On(t2.C("ItemId").Eq(t4.C("Id")))
				return NewSelector[Order](db).From(t5)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM " +
					"((`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id`=`t2`.`order_id`) " +
					"JOIN `item` AS `t4` ON `t2`.`item_id`=`t4`.`id`);",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.s.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
