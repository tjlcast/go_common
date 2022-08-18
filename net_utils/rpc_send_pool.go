package net_utils

import (
	"net/rpc"
	"sync"
)

var tcpClientPool = make(map[string]*rpc.Client)
var mutex sync.Mutex

func getTcpClient(endpoint string) (*rpc.Client, error) {
	mutex.Lock()
	defer mutex.Unlock()

	client, ok := tcpClientPool[endpoint]
	if !ok {
		// 连接远程 rpc
		newClient, err := rpc.DialHTTP("tcp", endpoint)
		if err != nil {
			return nil, err
		}
		tcpClientPool[endpoint] = newClient
		return newClient, nil
	} else {
		return client, nil
	}
}

func removeTcpClient(endpoint string) {
	mutex.Lock()
	defer mutex.Unlock()

	client, ok := tcpClientPool[endpoint]
	if ok {
		_ = client.Close()
		delete(tcpClientPool, endpoint)
	}
}

func SendTcp(endpoint string, serviceMethod string, args interface{}, reply interface{}) error {
	client, e := getTcpClient(endpoint)
	if e != nil {
		return e
	}

	// Call 调用命名函数method，等待其完成，并返回对应的错误状态
	e = client.Call(serviceMethod, args, reply)
	if e != nil {
		removeTcpClient(endpoint)
	}
	return e
}