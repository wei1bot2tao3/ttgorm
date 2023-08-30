package model

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"ttgorm/orm/internal/errs"
)

//	func Test_parseModel1(t *testing.T) {
//		testCass := []struct {
//			name      string
//			val       any
//			entity    any
//			wantModel *Model
//
//			wantErr error
//			opts    []Option
//		}{
//			{
//				name:   "test Model",
//				entity: TestModel{},
//
//				wantErr: errs.ErrPointerOnly,
//			},
//			{
//				name:   "test Model",
//				entity: &TestModel{},
//				wantModel: &Model{
//					TableName: "test_model",
//
//					Fields: []*Field{
//						{
//							ColName: "id",
//							GoName:  "Id",
//							Type:    reflect.TypeOf(int64(0)),
//							Offset:  0,
//						},
//						{
//							ColName: "first_name",
//							GoName:  "FirstName",
//							Type:    reflect.TypeOf(""),
//							Offset:  8,
//						},
//						{
//							ColName: "last_name",
//							GoName:  "LastName",
//							Type:    reflect.TypeOf(&sql.NullString{}),
//							Offset:  32,
//						},
//						{
//							ColName: "age",
//							GoName:  "Age",
//							Type:    reflect.TypeOf(int8(0)),
//							Offset:  24,
//						},
//					},
//				},
//			},
//		}
//		r := &registry{}
//		for _, tc := range testCass {
//			t.Run(tc.name, func(t *testing.T) {
//				res, err := r.Registry(tc.entity)
//				assert.Equal(t, tc.wantErr, err)
//				if err != nil {
//					return
//				}
//				filedMap := make(map[string]*Field)
//				columnMap := make(map[string]*Field)
//
//				for _, f := range tc.wantModel.Fields {
//					filedMap[f.GoName] = f
//					columnMap[f.ColName] = f
//
//				}
//				fmt.Println(filedMap)
//				tc.wantModel.FieldsMap = filedMap
//				tc.wantModel.ColumnMap = columnMap
//				assert.Equal(t, tc.wantModel, res)
//			})
//		}
//	}
//
//	func TestRegistry_get1(t *testing.T) {
//		testCass := []struct {
//			name string
//
//			entity    any
//			wantErr   error
//			wantModel *Model
//			wantRes   map[string]any
//			cacheSize int
//			opts      []Option
//		}{
//			{
//				name:   "test Model",
//				entity: &TestModel{},
//				wantModel: &Model{
//					TableName: "test_model",
//					Fields: []*Field{
//						{
//							ColName: "id",
//							GoName:  "Id",
//							Type:    reflect.TypeOf(int64(0)),
//							Offset:  0,
//						},
//						{
//							ColName: "first_name",
//							GoName:  "FirstName",
//							Type:    reflect.TypeOf(""),
//							Offset:  8,
//						},
//						{
//							ColName: "last_name",
//							GoName:  "LastName",
//							Type:    reflect.TypeOf(&sql.NullString{}),
//							Offset:  32,
//						},
//						{
//							ColName: "age",
//							GoName:  "Age",
//							Type:    reflect.TypeOf(int8(0)),
//							Offset:  24,
//						},
//					},
//				},
//
//				cacheSize: 1,
//			},
//
//			{
//				name: "tag",
//				entity: func() any {
//					type TagTable struct {
//						FirstName string `orm:"column=first_name_t"`
//					}
//					return &TagTable{}
//				}(),
//
//				wantModel: &Model{
//					TableName: "tag_table",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name_t",
//						},
//					},
//				},
//			},
//
//			{
//				name: "empty column",
//				entity: func() any {
//					type TagTable struct {
//						FirstName string
//					}
//					return &TagTable{}
//				}(),
//
//				wantModel: &Model{
//					TableName: "tag_table",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name",
//						},
//					},
//				},
//			},
//
//			{
//				name: "empty column",
//				entity: func() any {
//					type TagTable struct {
//						FirstName string `orm:"first_name"`
//					}
//					return &TagTable{}
//				}(),
//
//				wantModel: &Model{
//					TableName: "tag_table",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name",
//						},
//					},
//				},
//				wantErr: errs.NewErrInvalidTagContent("first_name"),
//			},
//			{
//				name:   "table Name",
//				entity: &CustomTableName{FirstName: "Tete1"},
//
//				wantModel: &Model{
//					TableName: "custom_table_nameTete1",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name",
//						},
//					},
//				},
//			},
//
//			{
//				name:   "table Name ptr",
//				entity: &CustomTableNamePtr{FirstName: "Tete1"},
//
//				wantModel: &Model{
//					TableName: "custom_table_nameTete1",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name",
//						},
//					},
//				},
//			},
//
//			{
//				name:   "table Name ptr empty",
//				entity: &CustomTableNameEmpty{FirstName: "Tete1"},
//
//				wantModel: &Model{
//					TableName: "custom_table_name_empty",
//					FieldsMap: map[string]*Field{
//						"FirstName": {
//							ColName: "first_name",
//						},
//					},
//				},
//			},
//		}
//
//		r := NewRegistry()
//		for _, tc := range testCass {
//			t.Run(tc.name, func(t *testing.T) {
//				res, err := r.Get(tc.entity)
//				assert.Equal(t, tc.wantErr, err)
//				if err != nil {
//					return
//				}
//				filedMap := make(map[string]*Field)
//				columnMap := make(map[string]*Field)
//
//				for _, f := range tc.wantModel.Fields {
//					filedMap[f.GoName] = f
//					columnMap[f.ColName] = f
//				}
//				fmt.Println(filedMap)
//				tc.wantModel.FieldsMap = filedMap
//				tc.wantModel.ColumnMap = columnMap
//				assert.Equal(t, tc.wantModel, res)
//			})
//		}
//	}
func Test_registry_Register(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:    "struct",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				Fields: []*Field{
					{
						ColName: "id",
						GoName:  "Id",
						Type:    reflect.TypeOf(int64(0)),
					},
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
						Offset:  8,
					},
					{
						ColName: "age",
						GoName:  "Age",
						Type:    reflect.TypeOf(int8(0)),
						Offset:  24,
					},
					{
						ColName: "last_name",
						GoName:  "LastName",
						Type:    reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
				},
			},
		},
		{
			name:    "map",
			entity:  map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "slice",
			entity:  []int{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "basic types",
			entity:  0,
			wantErr: errs.ErrPointerOnly,
		},
	}

	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Registry(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tc.wantModel.Fields {
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tc.wantModel.FieldsMap = fieldMap
			tc.wantModel.ColumnMap = columnMap
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name string

		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				Fields: []*Field{
					{
						ColName: "id",
						GoName:  "Id",
						Type:    reflect.TypeOf(int64(0)),
					},
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
						Offset:  8,
					},
					{
						ColName: "age",
						GoName:  "Age",
						Type:    reflect.TypeOf(int8(0)),
						Offset:  24,
					},
					{
						ColName: "last_name",
						GoName:  "LastName",
						Type:    reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
				},
			},
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
				Fields: []*Field{
					{
						ColName: "first_name_t",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				TableName: "tag_table",
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "column only",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			name: "ignore tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc=abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				TableName: "tag_table",
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				TableName: "custom_table_name",
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:   "table name ptr",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				TableName: "custom_table_name",
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Type:    reflect.TypeOf(""),
					},
				},
			},
		},
	}
	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			fieldMap := make(map[string]*Field)
			columnMap := make(map[string]*Field)
			for _, f := range tc.wantModel.Fields {
				fieldMap[f.GoName] = f
				columnMap[f.ColName] = f
			}
			tc.wantModel.FieldsMap = fieldMap
			tc.wantModel.ColumnMap = columnMap

			assert.Equal(t, tc.wantModel, m)

			typ := reflect.TypeOf(tc.entity)
			cache, ok := r.(*registry).models.Load(typ)
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
