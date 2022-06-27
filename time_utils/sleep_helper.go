package time_utils

import "time"

func WaitSeconds(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
}