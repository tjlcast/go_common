package arr_utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestArrString_Append(t *testing.T) {
	var arr = ArrString{"123", "456"}

	arr.Append("789")

	join := strings.Join(arr, "")

	if len(arr) != 3 {
		assert.Fail(t, "The length is not right.")
	}
	if join != "123456789" {
		assert.Fail(t, "The content is not right.")
	}
}
