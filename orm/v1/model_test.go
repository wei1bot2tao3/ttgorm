package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"ttgorm/orm/internal/errs"
)

func TestName(t *testing.T) {
	testCass := []struct {
		name string

		entity    any
		wantModel *model
		wantErr   error
	}{
		{
			name:   "test model",
			entity: TestModel{},
			//wantModel: &model{
			//	tableName: "test_model",
			//	fields: map[string]*field{
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
			name:   "test model",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
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
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, res)
		})
	}
}
