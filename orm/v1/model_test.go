package v1

import (
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
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
	}
	r := &registry{}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Registry(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, res)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCass := []struct {
		name string

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
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
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
				fields: map[string]*Field{
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
				fields: map[string]*Field{
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
				fields: map[string]*Field{
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
				fields: map[string]*Field{
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
				fields: map[string]*Field{
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
				fields: map[string]*Field{
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
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
			typ := reflect.TypeOf(tc.entity)
			cache, ok := r.models.Load(typ)

			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, cache)
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
			fd, ok := res.fields[tc.filed]
			require.True(t, ok)
			assert.Equal(t, tc.wantRes, fd.colName)

		})
	}
}
