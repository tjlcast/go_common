package mock_utils

import (
	"fmt"
	"github.com/tjlcast/go_common/executor_utils"
	"github.com/tjlcast/go_common/log_utils"
	"github.com/tjlcast/go_common/time_utils"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// @ http://c.biancheng.net/view/5124.html
var checkArgReg = regexp.MustCompile("\\d+!\\d+")

func PQWrap(inFunc func(id, arg string) error) func(id, arg string) {
	return func(id, arg string) {
		// check arg
		s := checkArgReg.FindString(arg)
		if s == "" {
			panic("Fail to check arg and should be xx!x")
		}

		split := strings.Split(arg, OP_MULIT_SPLITER)
		N, _ := strconv.Atoi(split[0])
		C, _ := strconv.Atoi(split[1])

		log_utils.TestLog(fmt.Sprintf("PQ test => N: %d C: %d.\n", N, C))

		var sCount uint64 = 0
		var fCount uint64 = 0
		var qCost int64 = 0

		pqPool, _ := executor_utils.NewPool(uint64(C))

		var wg sync.WaitGroup
		wg.Add(N)

		totalCost := time_utils.MarkTime()
		for i := 0; i < N; i++ {
			no := i
			pqTask := &executor_utils.Task{}
			pqTask.Handler = func(v ...interface{}) {
				qTime := time_utils.MarkTime()
				e := inFunc(strconv.Itoa(no), arg)
				if e != nil {
					atomic.AddUint64(&fCount, 1)
				} else {
					atomic.AddUint64(&sCount, 1)
				}
				qgap := qTime.Gap()
				atomic.AddInt64(&qCost, int64(qgap))
				wg.Done()
			}
			err := pqPool.Put(pqTask)
			if err != nil {
				log_utils.TestLog(fmt.Sprintf("Fail to put task(%d) and error: %s.\n", no, err.Error()))
			}
		}
		wg.Wait()
		log_utils.TestLog(fmt.Sprintf("TotalCost: %d countCost: %d ssize: %d fsize: %d.\n", totalCost.Gap(), qCost, sCount, fCount))
	}
}

