package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCasses := []struct {
		name string

		builder    QueryBuilder
		wantQuerry *Query
		wantErr    error
	}{
		{
			name:    "select no form",
			builder: &Selector[TestModel]{},
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "form",
			builder: (&Selector[TestModel]{}).Form("test_model1"),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_model1`;",
				Args: nil,
			},
		},
		{
			name:    "empty form ",
			builder: (&Selector[TestModel]{}).Form(""),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel`;",
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
			builder: (&Selector[TestModel]{}).Where(C("Id").Eq(18)),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel` WHERE `Id`=?;",
				Args: []any{18},
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(C("Id").Eq(18).And(C("Id").Eq(11))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel` WHERE (`Id`=?)AND(`Id`=?);",
				Args: []any{18, 11},
			},
		},
		{
			name:    "not",
			builder: (&Selector[TestModel]{}).Where(Not(C("Id").Eq(18))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel` WHERE NOT(`Id`=?);",
				Args: []any{18},
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(C("Id").Eq(18).Or(C("Id").Eq(11))),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel` WHERE (`Id`=?)OR(`Id`=?);",
				Args: []any{18, 11},
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `TestModel`;",
				Args: nil,
			},
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
