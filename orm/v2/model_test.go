package v1

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"ttgorm/orm/internal/errs"
)

func TestName(t *testing.T) {
	testCass := []struct {
		name      string
		val       any
		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "test Model",
			entity: TestModel{},
			//wantModel: &Model{
			//	tableName: "test_model",
			//	fields: map[string]*Field{
			//		"Id": {
			//			colName: "id",
			//		},
			//		"FirstName": {
			//			colName: "first_name",
			//		},
			//		"LastName": {
			//			colName: "last_name",
			//		},
			//		"Age": {
			//			colName: "age",
			//		},
			//	},
			//},
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
			res, err := r.Registry(tc.entity)
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
