package reflect

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestInterateFields(t *testing.T) {
	testCass := []struct {
		name   string
		entity any

		wantErr error
		wantRes map[string]any
	}{
		{
			name: "struct",
			entity: User{
				Name: "TT",
				age:  23,
			},
			wantRes: map[string]any{
				"Name": "TT",
				"age":  0,
			},
		},
		{
			name: "struct",
			entity: &User{
				Name: "TT",
				age:  23,
			},
			wantRes: map[string]any{
				"Name": "TT",
				"age":  0,
			},
		},
		{
			name:    "struct",
			entity:  19,
			wantErr: errors.New("不支持类型"),
		},

		{
			name: "struct",
			entity: func() **User {
				res := &User{
					Name: "TT",
					age:  22,
				}
				return &res
			}(),
			wantRes: map[string]any{
				"Name": "TT",
				"age":  0,
			},
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("不支持nil"),
		},
		{
			name:    "nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := InterateFields(tc.entity)

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type User struct {
	Name string
	age  int
}

func TestSetField(t *testing.T) {
	testCass := []struct {
		name       string
		entity     any
		field      string
		newVal     any
		wantErr    error
		wantRes    map[string]any
		wantEntity any
	}{
		{
			name: "struct",
			entity: User{
				Name: "tom",
			},
			field:  "Name",
			newVal: "Jerry",
			wantEntity: User{
				Name: "Jerry",
			},
			wantErr: errors.New("不允许修改"),
		},
		{
			name: "pointer",
			entity: &User{
				Name: "tom",
			},
			field:  "Name",
			newVal: "Jerry",
			wantEntity: &User{
				Name: "Jerry",
			},
		},
		{
			name: "pointer exported",
			entity: &User{
				Name: "tom",
				age:  19,
			},
			field:  "age",
			newVal: "12",
			wantEntity: &User{
				Name: "Jerry",
				age:  19,
			},
			wantErr: errors.New("不允许修改"),
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}

	var i = 0
	prt := &i
	reflect.ValueOf(prt).Elem().Set(reflect.ValueOf(12))
	assert.Equal(t, 12, i)
}

func TestFr(t *testing.T) {
	type Usaer struct {
		Name string
		age  int
		Sex  string
	}
	ua := Usaer{
		Name: "撒",
		age:  12,
		Sex:  "男",
	}

	typ := reflect.TypeOf(&ua)

	val := reflect.ValueOf(ua)

	fmt.Println(typ.Kind())
	fmt.Println(val.Kind())
	typ = typ.Elem()
	numfidld := typ.NumField()
	fmt.Println(numfidld)
	//res:=make(map[string]any,numfidld)
	for i := 0; i < numfidld; i++ {
		fileVal := val.Field(i)
		filetype := typ.Field(i)
		if filetype.IsExported() {

		} else {
			fmt.Println(filetype.Name, fileVal)
		}
	}

}
