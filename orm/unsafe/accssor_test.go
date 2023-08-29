package unsafe

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnSafeAccessor_Filed(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	accessor := NewUnSafeAccessor(&User{
		Name: "Tom",
		Age:  18,
	})

	val, err := accessor.Filed("Name")
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, "Tom", val)

	err = accessor.Set("Name", "Jerry")
	assert.NoError(t, err)
	val, err = accessor.Filed("Name")
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, "Jerry", val)

}
