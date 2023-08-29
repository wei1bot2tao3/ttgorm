package unsafe

import "testing"

func TestPrintFieldOffset(t *testing.T) {
	testCass := []struct {
		name   string
		entity any

		wantErr error
		wantRes map[string]any
	}{
		{
			name:   "user",
			entity: User{},
		},
	}
	for _, tc := range testCass {
		t.Run(tc.name, func(t *testing.T) {
			PrintFieldOffset(tc.entity)
		})
	}
}

type User struct {
	// 0
	Name string
	//16
	Age int32
	//20   为什么？
	Agev1 int32
	//24
	Alias []string
	//48
	Addresss string
}
