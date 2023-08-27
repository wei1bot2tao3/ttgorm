package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArry(t *testing.T) {

	testCass := []struct {
		name string

		entity  any
		wantErr error
		wantRes []any
	}{
		{
			name:    "[]int",
			entity:  [3]int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
		{
			name:    "slice",
			entity:  []int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
		{
			name:    "map",
			entity:  []int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			vales, err := IterateSlice(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, vales)
		})
	}
}

func TestIterateMap(t *testing.T) {
	testCass := []struct {
		name string

		entity    any
		wantErr   error
		wantKeys  []any
		wantValue []any
	}{
		{
			name: "map",
			entity: map[string]string{
				"A": "a",
				"B": "b",
			},
			wantKeys:  []any{"A", "B"},
			wantValue: []any{"a", "b"},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			key, value, err := IterateMap(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.EqualValues(t, tc.wantKeys, key)
			assert.EqualValues(t, tc.wantValue, value)

		})
	}
}