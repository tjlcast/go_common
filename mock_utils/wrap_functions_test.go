package mock_utils

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestPQWrapWrongArg(t *testing.T) {
	var count int64

	wrap := PQWrap(func(id, arg string) error {
		//fmt.Println(fmt.Sprintf("Id %s, Args: %s\n", id, arg))
		atomic.AddInt64(&count, 1)
		return nil
	})

	defer func() {
		if err:= recover(); err==nil {
			assert.Fail(t, "Should panic with wrong args")
		}
	}()
	wrap("000", "100-1")
}

func TestPQWrap100N1C(t *testing.T) {
	var count int64

	wrap := PQWrap(func(id, arg string) error {
		//fmt.Println(fmt.Sprintf("Id %s, Args: %s\n", id, arg))
		atomic.AddInt64(&count, 1)
		return nil
	})

	wrap("000", "100!1")

	if count != 100 {
		assert.Fail(t, "The not do 100 times.")
	}
}

func TestPQWrap100N4C(t *testing.T) {
	var count int64

	wrap := PQWrap(func(id, arg string) error {
		//fmt.Println(fmt.Sprintf("Id %s, Args: %s\n", id, arg))
		atomic.AddInt64(&count, 1)
		return nil
	})

	wrap("000", "100!4")

	if count != 100 {
		assert.Fail(t, "The not do 100 times.")
	}
}
