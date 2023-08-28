package v1

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"ttgorm/orm/internal/errs"
)

func TestSelector_Build(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)
	testCasses := []struct {
		name string

		builder    QueryBuilder
		wantQuerry *Query
		wantErr    error
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
