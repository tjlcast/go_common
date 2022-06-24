package net_utils

import (
	"net"
	"sync"
)

var Listener net.Listener = nil

var listerMapping = make(map[string]net.Listener)
var locker = sync.Mutex{}

func GetOrBuildListener(endpoint string) (net.Listener, error) {
	locker.Lock()
	defer func() {
		locker.Unlock()
	}()

	var err error
	lister, ok := listerMapping[endpoint]
	if !ok {
		lister, err = net.Listen("tcp", endpoint)
		if err != nil {
			return nil, err
		}
		listerMapping[endpoint] = lister
	}
	return lister, nil
}


