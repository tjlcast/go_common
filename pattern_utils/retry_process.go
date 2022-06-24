package pattern_utils

import (
	"fmt"
	"log"
	"time"
)

// Try provides kissrata's retry strategy.
// It calls the given function up to five times or until the given function
// returns nil, whichever comes first.
func Try(f func() error, header string) error {
	if header != "" {
		Log(header)
	}
	err := f()
	for tryNum := 2; tryNum <= 5 && err != nil; tryNum++ {
		Log(fmt.Sprintf("%s", err))
		time.Sleep(200 * time.Millisecond)
		Log(fmt.Sprintf("(try number %d) "+header, tryNum))
		err = f()
	}
	return err
}


//Log prints logging info
func Log(msg string) {
	log.Println(msg)
}