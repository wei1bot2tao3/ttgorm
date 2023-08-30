package model

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
		fields    []*Fields
		wantErr   error
		opts      []Option
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
				TableName: "test_model",
			},
			fields: []*Fields{
				{
					ColName: "id",
					GoName:  "Id",
					Type:    reflect.TypeOf(int64(0)),
					Offset:  0,
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
					Offset:  8,
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Type:    reflect.TypeOf(&sql.NullString{}),
					Offset:  32,
				},
				{
					ColName: "age",
					GoName:  "Age",
					Type:    reflect.TypeOf(int8(0)),
					Offset:  24,
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
			filedMap := make(map[string]*Fields)
			columnMap := make(map[string]*Fields)

			for _, f := range tc.fields {
				filedMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			fmt.Println(filedMap)
			tc.wantModel.FieldsMap = filedMap
			tc.wantModel.ColumnMap = columnMap
			assert.Equal(t, tc.wantModel, res)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCass := []struct {
		name      string
		fields    []*Fields
		entity    any
		wantErr   error
		wantModel *Model
		wantRes   map[string]any
		cacheSize int
		opts      []Option
	}{
		{
			name:   "test Model",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
			},
			fields: []*Fields{
				{
					ColName: "id",
					GoName:  "Id",
					Type:    reflect.TypeOf(int64(0)),
					Offset:  0,
				},
				{
					ColName: "first_name",
					GoName:  "FirstName",
					Type:    reflect.TypeOf(""),
					Offset:  8,
				},
				{
					ColName: "last_name",
					GoName:  "LastName",
					Type:    reflect.TypeOf(&sql.NullString{}),
					Offset:  32,
				},
				{
					ColName: "age",
					GoName:  "Age",
					Type:    reflect.TypeOf(int8(0)),
					Offset:  24,
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
				TableName: "tag_table",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name_t",
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
				TableName: "tag_table",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name",
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
				TableName: "tag_table",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
			wantErr: errs.NewErrInvalidTagContent("first_name"),
		},
		{
			name:   "table Name",
			entity: &CustomTableName{FirstName: "Tete1"},

			wantModel: &Model{
				TableName: "custom_table_nameTete1",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},

		{
			name:   "table Name ptr",
			entity: &CustomTableNamePtr{FirstName: "Tete1"},

			wantModel: &Model{
				TableName: "custom_table_nameTete1",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},

		{
			name:   "table Name ptr empty",
			entity: &CustomTableNameEmpty{FirstName: "Tete1"},

			wantModel: &Model{
				TableName: "custom_table_name_empty",
				FieldsMap: map[string]*Fields{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},
	}

	r := NewRegistry()
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			filedMap := make(map[string]*Fields)
			columnMap := make(map[string]*Fields)

			for _, f := range tc.fields {
				filedMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			fmt.Println(filedMap)
			tc.wantModel.FieldsMap = filedMap
			tc.wantModel.ColumnMap = columnMap
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
	r := NewRegistry()
	m, err := r.Registry(&TestModel{}, WithTableName("test_model_ttt"))
	require.NoError(t, err)
	assert.Equal(t, "test_model_ttt", m.TableName)
}

func TestModleWithColumnName(t *testing.T) {
	testCass := []struct {
		name    string
		filed   string
		colName string
		wantErr error
		wantRes string
		opts    []Option
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
			r := NewRegistry()
			res, err := r.Registry(&TestModel{}, WithColumnName(tc.filed, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := res.FieldsMap[tc.filed]
			require.True(t, ok)
			assert.Equal(t, tc.wantRes, fd.ColName)

		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
