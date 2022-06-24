package chan_utils

import (
	"errors"
	"time"
)

func ReadWithSelect(ch chan string) (x string, err error) {
	timeout := time.NewTimer(time.Microsecond * 5000)

	select {
	case x = <-ch:
		return x, nil
	case <-timeout.C:
		return "", errors.New("read time out")
	}
}

func WriteChWithSelect(ch chan interface{}, msg interface{}) error {
	timeout := time.NewTimer(time.Microsecond * 5000)

	select {
	case ch <- msg:
		return nil
	case <-timeout.C:
		return errors.New("write time out")
	}
}
