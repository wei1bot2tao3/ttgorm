package reflect

import (
	"github.com/stretchr/testify/assert"
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
			name:   "struct",
			entity: 19,
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
