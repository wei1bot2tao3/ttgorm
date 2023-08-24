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
		{
			name:    "form db ",
			builder: (&Selector[TestModel]{}).Form("test_db.test_model"),
			wantQuerry: &Query{
				SQL:  "SELECT * FORM `test_db`.`test_model`;",
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
