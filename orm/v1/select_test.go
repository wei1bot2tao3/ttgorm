package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"ttgorm/orm/internal/errs"
)

func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)

	testCasses := []struct {
		name string

		builder    QueryBuilder
		wantQuerry *Query
		wantErr    error
		wantRes    *TestModel
	}{
		{
			name:    "select no form",
			builder: NewSelector[TestModel](db),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "form",
			builder: NewSelector[TestModel](db).Form("test_model"),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty form ",
			builder: NewSelector[TestModel](db).Form(""),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model`;",
				Args: nil,
			},
		},
		//{
		//	name:    "form db ",
		//	builder: (&Selector[TestModel]{}).Form("test_db.test_model"),
		//	wantQuerry: &Query{
		//		SQL:  "SELECT * FORM `test_db`.`test_model`;",
		//		Args: nil,
		//	},
		//},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18)),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model` WHERE `id`=?;",
				Args: []any{18},
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18).And(C("Id").Eq(11))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model` WHERE (`id`=?)AND(`id`=?);",
				Args: []any{18, 11},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(C("Id").Eq(18))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model` WHERE NOT(`id`=?);",
				Args: []any{18},
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(C("Id").Eq(18).Or(C("Id").Eq(11))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model` WHERE (`id`=?)OR(`id`=?);",
				Args: []any{18, 11},
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(Not(C("jkd").Eq(18))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model`;",
				Args: nil,
			},
			wantErr: errs.NewErrUnknownField("jkd"),
		},
	}

	for _, tc := range testCasses {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuerry, q)
		})
	}
}

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

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
		{
			name: "scan error",
			s:    NewSelector[TestModel](db).Where(C("Id").LT(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
			wantErr: errs.ErrNoRows,
		},
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
