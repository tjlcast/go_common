package log_utils

import (
	"fmt"
	"github.com/tjlcast/go_common/rand_utils"
	"strconv"
	"sync"
	"testing"
)

type Recorder struct {
	rec map[string]int
	wg sync.RWMutex
}

func NewRecorder() *Recorder {
	return &Recorder{
		rec: make(map[string]int),
	}
}

func (r Recorder) Incr(key string) int {
	r.wg.Lock()
	defer r.wg.Unlock()

	v, ok := r.rec[key]
	if !ok {
		v = 0
	}

	v += 1
	r.rec[key] = v

	return v
}

func TestRandLog(t *testing.T) {
	NUM := 10000

	recorder := NewRecorder()

	for i:=0; i<NUM; i++ {
		randInt := rand_utils.GenRandInt(10)
		incr := recorder.Incr(strconv.Itoa(randInt))
		fmt.Println(incr)
	}
}
