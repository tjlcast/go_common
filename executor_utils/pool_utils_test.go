package executor_utils

import (
	"fmt"
	"github.com/tjlcast/go_common/time_utils"
	"strconv"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool, _ := NewPool(2)
	tasks := make([]*Task, 10)

	for j:=0; j<10; j++ {
		i := j
		go func() {
			task := &Task{}
			tasks[i] = task
			task.Params = []interface{}{strconv.Itoa(i)}
			task.Handler = func(v ...interface{}) {
				for !task.Interupt {
					fmt.Println(v[0].(string) + " is running...")
					time_utils.WaitSeconds(1)
				}
			}
			err := pool.Put(task)
			if err != nil {
				fmt.Printf("Fail to put a task.")
			}
		}()
	}

	closedNo := 0
	for {
		select {
		case <- time.After(3 * time.Second):
			tasks[closedNo].Interupt = true
			fmt.Println("Close " + strconv.Itoa(closedNo))
			closedNo += 1
			if closedNo >= 10 {
				return
			}
		}
	}
	pool.Close()
}

func TestNewPool1(t *testing.T) {
	pool, _ := NewPool(2)
	tasks := make([]*Task, 10)

	for j:=0; j<10; j++ {
		i := j

		task := &Task{}
		tasks[i] = task
		fmt.Println("Submit the task: " + strconv.Itoa(i))
		task.Params = []interface{}{strconv.Itoa(i)}
		task.Handler = func(v ...interface{}) {
			for !task.Interupt {
				fmt.Println(v[0].(string) + " is running...")
				time_utils.WaitSeconds(1)
			}
		}
		err := pool.Put(task)
		if err != nil {
			fmt.Printf("Fail to put a task.")
		}
	}
}
