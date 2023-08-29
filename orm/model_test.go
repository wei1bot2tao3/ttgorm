package orm

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"ttgorm/orm/internal/errs"
)

func Test_parseModel(t *testing.T) {
	testCass := []struct {
		name      string
		val       any
		entity    any
		wantModel *Model
		fields    []*Field
		wantErr   error
		opts      []ModelOption
	}{
		{
			name:   "test Model",
			entity: TestModel{},

			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "test Model",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
			},
			fields: []*Field{
				{
					colName: "id",
					GOName:  "Id",
					typ:     reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					GOName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					GOName:  "LastName",
					typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					GOName:  "Age",
					typ:     reflect.TypeOf(int8(0)),
				},
			},
		},
	}
	r := &registry{}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Registry(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			filedMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)

			for _, f := range tc.fields {
				filedMap[f.GOName] = f
				columnMap[f.colName] = f
			}
			fmt.Println(filedMap)
			tc.wantModel.fieldsMap = filedMap
			tc.wantModel.columnMap = columnMap
			assert.Equal(t, tc.wantModel, res)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCass := []struct {
		name      string
		fields    []*Field
		entity    any
		wantErr   error
		wantModel *Model
		wantRes   map[string]any
		cacheSize int
		opts      []ModelOption
	}{
		{
			name:   "test Model",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
			},
			fields: []*Field{
				{
					colName: "id",
					GOName:  "Id",
					typ:     reflect.TypeOf(int64(0)),
				},
				{
					colName: "first_name",
					GOName:  "FirstName",
					typ:     reflect.TypeOf(""),
				},
				{
					colName: "last_name",
					GOName:  "LastName",
					typ:     reflect.TypeOf(&sql.NullString{}),
				},
				{
					colName: "age",
					GOName:  "Age",
					typ:     reflect.TypeOf(int8(0)),
				},
			},
			cacheSize: 1,
		},

		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),

			wantModel: &Model{
				tableName: "tag_table",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
		},

		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string
				}
				return &TagTable{}
			}(),

			wantModel: &Model{
				tableName: "tag_table",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"first_name"`
				}
				return &TagTable{}
			}(),

			wantModel: &Model{
				tableName: "tag_table",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr: errs.NewErrInvalidTagContent("first_name"),
		},
		{
			name:   "table Name",
			entity: &CustomTableName{FirstName: "Tete1"},

			wantModel: &Model{
				tableName: "custom_table_nameTete1",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name:   "table Name ptr",
			entity: &CustomTableNamePtr{FirstName: "Tete1"},

			wantModel: &Model{
				tableName: "custom_table_nameTete1",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name:   "table Name ptr empty",
			entity: &CustomTableNameEmpty{FirstName: "Tete1"},

			wantModel: &Model{
				tableName: "custom_table_name_empty",
				fieldsMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
	}

	r := newRegistry()
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			filedMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)

			for _, f := range tc.fields {
				filedMap[f.GOName] = f
				columnMap[f.colName] = f
			}
			fmt.Println(filedMap)
			tc.wantModel.fieldsMap = filedMap
			tc.wantModel.columnMap = columnMap
			assert.Equal(t, tc.wantModel, res)
		})
	}
}

type CustomTableName struct {
	FirstName string
}

func (c *CustomTableName) TableName() string {
	tabl := "custom_table_name" + c.FirstName
	return tabl
}

type CustomTableNamePtr struct {
	FirstName string
}

func (c *CustomTableNamePtr) TableName() string {
	tabl := "custom_table_name" + c.FirstName
	return tabl
}

type CustomTableNameEmpty struct {
	FirstName string
}

func (c *CustomTableNameEmpty) TableName() string {
	tabl := ""
	return tabl
}

func TestModelWithTableName(t *testing.T) {
	r := newRegistry()
	m, err := r.Registry(&TestModel{}, ModelWithTableName("test_model_ttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_ttt", m.tableName)
}

func TestModleWithColumnName(t *testing.T) {
	testCass := []struct {
		name    string
		filed   string
		colName string
		wantErr error
		wantRes string
		opts    []ModelOption
	}{
		{
			name:    "column came ",
			filed:   "FirstName",
			colName: "first_name_ccc",
			wantRes: "first_name_ccc",
		},
		{
			name:    "invalid column came ",
			filed:   "XXXXX",
			colName: "first_name_ccc",
			wantRes: "first_name_ccc",
			wantErr: errs.NewErrUnknownField("XXXXX"),
		},
	}

	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			r := newRegistry()
			res, err := r.Registry(&TestModel{}, ModleWithColumnName(tc.filed, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := res.fieldsMap[tc.filed]
			require.True(t, ok)
			assert.Equal(t, tc.wantRes, fd.colName)

		})
	}
}
