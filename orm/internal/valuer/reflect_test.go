package valuer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"ttgorm/orm/model"
)

func Test_reflectValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewReflectValue)
	testSetColumns(t, NewUnsafeValue)
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func testSetColumns(t *testing.T, creator Creator) {
	testCass := []struct {
		name string

		entity  any
		wantErr error
		rows    func() *sqlmock.Rows
		wantRes map[string]any

		wantEntity any
	}{
		{
			name:   "set column",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", "18", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		// 不同位置
		{
			name:   "set column",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", "18", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
		{
			// 测试列的不同顺序
			name:   "order",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})
				rows.AddRow("1", "Jerry", "Tom", "18")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试列的不同顺序
			name:   "partial columns",
			entity: &TestModel{},
			rows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name"})
				rows.AddRow("1", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}
	r := model.NewRegistry()
	mockDb, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDb.Close()
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {

			mockRows := tc.rows()
			mock.ExpectQuery("SELECT XXX").WillReturnRows(mockRows)
			rows, err := mockDb.Query("SELECT XXX")
			require.NoError(t, err)
			rows.Next()

			model, err := r.Get(tc.entity)
			require.NoError(t, err)
			val := creator(model, tc.entity)

			err = val.SetColumns(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}
