package additional

type MyStruct struct {
	// 第一部分是用户必须输入的字段
	id   uint64
	name string

	// 第二部分是可选的字段
	address string
	// 可以很多字段
	filed1 int
	filed2 int
}

type MyStructOption func(myStruct *MyStruct)

type MyStructOptionErr func(myStruct *MyStruct) error

func WithAddress(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.address = address
	}
}

func WithField1andField2(field1 int, filed2 int) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.filed1 = field1
		myStruct.filed2 = filed2
	}
}

// NewMyStruct 都是私有的所以会 new一个
func NewMyStruct(id uint64, name string, opts ...MyStructOption) *MyStruct {
	res := &MyStruct{
		id:   id,
		name: name,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

// NewMyStruct 都是私有的所以会 new一个
func NewMyStructErr(id uint64, name string, opts ...MyStructOptionErr) (*MyStruct, error) {
	res := &MyStruct{
		id:   id,
		name: name,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
