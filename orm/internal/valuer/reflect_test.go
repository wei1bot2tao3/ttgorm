package valuer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	v1 "ttgorm/orm/model"
)

func Test_reflectValue_SetColumns(t *testing.T) {
	testCass := []struct {
		name      string
		etity     any
		wantErr   error
		wantRes   map[string]any
		rows      func() *sqlmock.Rows
		wantEtity any
	}{
		{
			name:  "set colums",
			etity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", "18", "Jerry")
				return rows
			},
			wantEtity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}
	r := v1.NewRegistry()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			model, err := r.Get(tc.etity)
			require.NoError(t, err)

			val := NewReflectValue(model, tc.etity)

			mockRows := tc.rows()
			mock.ExpectQuery("SELECT XX").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT XX")
			require.NoError(t, err)
			rows.Next()
			err = val.SetColumns(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEtity, tc.etity)

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
