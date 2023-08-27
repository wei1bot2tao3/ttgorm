package reflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	types "ttgorm/orm/reflect/type"
)

func TestNteraFunc(t *testing.T) {
	testCass := []struct {
		name string

		entiy any

		wantErr error
		wantRes map[string]FuncInfo
	}{
		{
			name:  "struct",
			entiy: types.NewUser("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:       "GetAge",
					InputTypes: []reflect.Type{reflect.TypeOf(types.User{})},
					OutTypes:   []reflect.Type{reflect.TypeOf(0)},
					Result:     []any{18},
				},

				//"ChangeName": {
				//	Name:       "ChangeName",
				//	InputTypes: []reflect.Type{reflect.TypeOf(types.User{})},
				//	OutTypes:   []reflect.Type{reflect.TypeOf(0)},
				//	Result:     []any{18},
				//},
			},
		},
		{
			name:  "pointer",
			entiy: types.NewUserPtr("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:       "GetAge",
					InputTypes: []reflect.Type{reflect.TypeOf(&types.User{})},
					OutTypes:   []reflect.Type{reflect.TypeOf(0)},
					Result:     []any{18},
				},

				"ChangeName": {
					Name:       "ChangeName",
					InputTypes: []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf(" ")},
					OutTypes:   []reflect.Type{},
					Result:     []any{},
				},
			},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res := IterateFunc(tc.entiy)
			//assert.Equal(t, tc.wantErr, err)
			//if err != nil {
			//	return
			//}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestIterateFunc(t *testing.T) {
	entiy := types.NewUser("Tom", 11)
	IterateFunc(entiy)

}
