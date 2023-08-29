package model

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"reflect"
//	"testing"
//	"ttgorm/orm/internal/errs"
//	"ttgorm/orm/v1"
//)
//
//func Test_parseModel(t *testing.T) {
//	testCass := []struct {
//		name      string
//		val       any
//		entity    any
//		wantModel *Model
//		fields    []*Field
//		wantErr   error
//		opts      []ModelOption
//	}{
//		{
//			name:   "test Model",
//			entity: v1.TestModel{},
//
//			wantErr: errs.ErrPointerOnly,
//		},
//		{
//			name:   "test Model",
//			entity: &v1.TestModel{},
//			wantModel: &Model{
//				TableName: "test_model",
//			},
//			fields: []*Field{
//				{
//					ColName: "id",
//					GOName:  "Id",
//					Typ:     reflect.TypeOf(int64(0)),
//				},
//				{
//					ColName: "first_name",
//					GOName:  "FirstName",
//					Typ:     reflect.TypeOf(""),
//				},
//				{
//					ColName: "last_name",
//					GOName:  "LastName",
//					Typ:     reflect.TypeOf(&sql.NullString{}),
//				},
//				{
//					ColName: "age",
//					GOName:  "Age",
//					Typ:     reflect.TypeOf(int8(0)),
//				},
//			},
//		},
//	}
//	r := &registry{}
//	for _, tc := range testCass {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := r.Registry(tc.entity)
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			filedMap := make(map[string]*Field)
//			columnMap := make(map[string]*Field)
//
//			for _, f := range tc.fields {
//				filedMap[f.GOName] = f
//				columnMap[f.ColName] = f
//			}
//			fmt.Println(filedMap)
//			tc.wantModel.FieldsMap = filedMap
//			tc.wantModel.ColumnMap = columnMap
//			assert.Equal(t, tc.wantModel, res)
//		})
//	}
//}
//
//func TestRegistry_get(t *testing.T) {
//	testCass := []struct {
//		name      string
//		fields    []*Field
//		entity    any
//		wantErr   error
//		wantModel *Model
//		wantRes   map[string]any
//		cacheSize int
//		opts      []ModelOption
//	}{
//		{
//			name:   "test Model",
//			entity: &v1.TestModel{},
//			wantModel: &Model{
//				TableName: "test_model",
//			},
//			fields: []*Field{
//				{
//					ColName: "id",
//					GOName:  "Id",
//					Typ:     reflect.TypeOf(int64(0)),
//				},
//				{
//					ColName: "first_name",
//					GOName:  "FirstName",
//					Typ:     reflect.TypeOf(""),
//				},
//				{
//					ColName: "last_name",
//					GOName:  "LastName",
//					Typ:     reflect.TypeOf(&sql.NullString{}),
//				},
//				{
//					ColName: "age",
//					GOName:  "Age",
//					Typ:     reflect.TypeOf(int8(0)),
//				},
//			},
//			cacheSize: 1,
//		},
//
//		{
//			name: "tag",
//			entity: func() any {
//				type TagTable struct {
//					FirstName string `orm:"column=first_name_t"`
//				}
//				return &TagTable{}
//			}(),
//
//			wantModel: &Model{
//				TableName: "tag_table",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name_t",
//					},
//				},
//			},
//		},
//
//		{
//			name: "empty column",
//			entity: func() any {
//				type TagTable struct {
//					FirstName string
//				}
//				return &TagTable{}
//			}(),
//
//			wantModel: &Model{
//				TableName: "tag_table",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name",
//					},
//				},
//			},
//		},
//
//		{
//			name: "empty column",
//			entity: func() any {
//				type TagTable struct {
//					FirstName string `orm:"first_name"`
//				}
//				return &TagTable{}
//			}(),
//
//			wantModel: &Model{
//				TableName: "tag_table",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name",
//					},
//				},
//			},
//			wantErr: errs.NewErrInvalidTagContent("first_name"),
//		},
//		{
//			name:   "table Name",
//			entity: &CustomTableName{FirstName: "Tete1"},
//
//			wantModel: &Model{
//				TableName: "custom_table_nameTete1",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name",
//					},
//				},
//			},
//		},
//
//		{
//			name:   "table Name ptr",
//			entity: &CustomTableNamePtr{FirstName: "Tete1"},
//
//			wantModel: &Model{
//				TableName: "custom_table_nameTete1",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name",
//					},
//				},
//			},
//		},
//
//		{
//			name:   "table Name ptr empty",
//			entity: &CustomTableNameEmpty{FirstName: "Tete1"},
//
//			wantModel: &Model{
//				TableName: "custom_table_name_empty",
//				FieldsMap: map[string]*Field{
//					"FirstName": {
//						ColName: "first_name",
//					},
//				},
//			},
//		},
//	}
//
//	r := NewRegistry()
//	for _, tc := range testCass {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := r.Get(tc.entity)
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			filedMap := make(map[string]*Field)
//			columnMap := make(map[string]*Field)
//
//			for _, f := range tc.fields {
//				filedMap[f.GOName] = f
//				columnMap[f.ColName] = f
//			}
//			fmt.Println(filedMap)
//			tc.wantModel.FieldsMap = filedMap
//			tc.wantModel.ColumnMap = columnMap
//			assert.Equal(t, tc.wantModel, res)
//		})
//	}
//}
//
//type CustomTableName struct {
//	FirstName string
//}
//
//func (c *CustomTableName) TableName() string {
//	tabl := "custom_table_name" + c.FirstName
//	return tabl
//}
//
//type CustomTableNamePtr struct {
//	FirstName string
//}
//
//func (c *CustomTableNamePtr) TableName() string {
//	tabl := "custom_table_name" + c.FirstName
//	return tabl
//}
//
//type CustomTableNameEmpty struct {
//	FirstName string
//}
//
//func (c *CustomTableNameEmpty) TableName() string {
//	tabl := ""
//	return tabl
//}
//
//func TestModelWithTableName(t *testing.T) {
//	r := NewRegistry()
//	m, err := r.Registry(&v1.TestModel{}, ModelWithTableName("test_model_ttt"))
//	require.NoError(t, err)
//	assert.Equal(t, "test_model_ttt", m.TableName)
//}
//
//func TestModleWithColumnName(t *testing.T) {
//	testCass := []struct {
//		name    string
//		filed   string
//		colName string
//		wantErr error
//		wantRes string
//		opts    []ModelOption
//	}{
//		{
//			name:    "column came ",
//			filed:   "FirstName",
//			colName: "first_name_ccc",
//			wantRes: "first_name_ccc",
//		},
//		{
//			name:    "invalid column came ",
//			filed:   "XXXXX",
//			colName: "first_name_ccc",
//			wantRes: "first_name_ccc",
//			wantErr: errs.NewErrUnknownField("XXXXX"),
//		},
//	}
//
//	for _, tc := range testCass {
//		t.Run(tc.name, func(t *testing.T) {
//			r := NewRegistry()
//			res, err := r.Registry(&v1.TestModel{}, ModleWithColumnName(tc.filed, tc.colName))
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			fd, ok := res.FieldsMap[tc.filed]
//			require.True(t, ok)
//			assert.Equal(t, tc.wantRes, fd.ColName)
//
//		})
//	}
//}
