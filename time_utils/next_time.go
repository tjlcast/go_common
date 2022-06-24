package time_utils

import (
	"sync"
	"time"
)

var lastSecond int64
var lastSecondMutex sync.Mutex

func GenNextSecondTs() int64 {
	lastSecondMutex.Lock()
	defer lastSecondMutex.Unlock()

	currentTs := time.Now().Unix()
	if currentTs == lastSecond {
		time.Sleep(time.Duration(1) * time.Second)
	}
	currentTs = time.Now().Unix()
	lastSecond = currentTs
	return lastSecond
}