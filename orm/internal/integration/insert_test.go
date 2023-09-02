package integration

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"ttgorm/orm"
	"ttgorm/orm/internal/test"
)

//func TestMySql(t *testing.T) {
//	testInsert(t, "mysql", "root:root@tcp(localhost:13306)/integration_test")
//
//}

//
//func testInsert(t *testing.T, driver string, DataSourcename string) {
//	db, err := orm.Open(driver, DataSourcename)
//	require.NoError(t, err)
//	testCass := []struct {
//		name string
//		i    *orm.Inserter[test.SimpleStruct]
//
//		wantErr      error
//		wantRes      map[string]any
//		wantAffected int64
//	}{
//		{
//			name:         "insert one",
//			i:            orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(14)),
//			wantAffected: 1,
//		},
//		{
//			name:         "insert multiple",
//			i:            orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(12), test.NewSimpleStruct(13)),
//			wantAffected: 2,
//		},
//		{
//			name:         "insert id",
//			i:            orm.NewInserter[test.SimpleStruct](db).Values(&test.SimpleStruct{Id: 19}),
//			wantAffected: 1,
//		},
//	}
//	for _, tc := range testCass {
//		t.Run(tc.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//			defer cancel()
//			res := tc.i.Exec(ctx)
//			affected, err := res.RowsAffected()
//			if err != nil {
//				fmt.Println(err)
//			}
//			require.NoError(t, res.Err())
//			assert.Equal(t, tc.wantAffected, affected)
//
//		})
//	}
//}

type InsertSuite struct {
	Suite
}

func TestMySQLInsert(t *testing.T) {
	suite.Run(t, &InsertSuite{
		Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (i *InsertSuite) TestInsert() {
	db := i.db
	t := i.T()
	testCases := []struct {
		name         string
		i            *orm.Inserter[test.SimpleStruct]
		wantAffected int64 // 插入行数
	}{
		{
			name:         "insert one",
			i:            orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(12)),
			wantAffected: 1,
		},
		{
			name: "insert multiple",
			i: orm.NewInserter[test.SimpleStruct](db).Values(
				test.NewSimpleStruct(13),
				test.NewSimpleStruct(14)),
			wantAffected: 2,
		},
		{
			name:         "insert id",
			i:            orm.NewInserter[test.SimpleStruct](db).Values(&test.SimpleStruct{Id: 15}),
			wantAffected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			res := tc.i.Exec(ctx)
			affected, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, tc.wantAffected, affected)
		})
	}
}

func (i *InsertSuite) SetupSuite() {
	db, err := orm.Open(i.driver, i.dsn)
	require.NoError(i.T(), err)
	i.db = db
}

type TtTest struct {
	Id   int
	Name string
}

func TestTtTable(t *testing.T) {
	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res := orm.NewInserter[TtTest](db).Values(&TtTest{Id: 2, Name: "TTs"}).Exec(ctx)
	if res.Err() != nil {
		fmt.Println(res.Err())
	}
	data, _ := orm.NewSelector[TtTest](db).Where(orm.C("Id").Eq(1)).Get(ctx)

	fmt.Println(data)

}
