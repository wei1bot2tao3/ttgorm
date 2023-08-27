package homework_delete

import (
	"testing"
	"ttgorm/orm/internal/errs"

	"github.com/stretchr/testify/assert"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no whereand form",
			builder: &Deleter[TestModel]{},
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "no where",
			builder: (&Deleter[TestModel]{}).From("`test_model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id`=?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id`=?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).Where(C("T").GT(12).Or(C("Id").LT(12))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`t`=?)OR(`id`<?);",
				Args: []any{12},
			},
			wantErr: errs.NewErrUnknownField("T"),
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
