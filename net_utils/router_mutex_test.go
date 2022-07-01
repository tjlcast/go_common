package net_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/rpc"
	"testing"
)

func TestNewRouterMulti(t *testing.T) {
	multi := NewRouterMulti(9999)

	multi.AddRestHandle(REST_GET, "/rest", RestApi)
	multi.AddRpcHandle(&RpcApi{})

	multi.Init()
	go multi.Loop()

	// 连接远程rpc
	rp, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	var bo *string
	err = rp.Call("RpcApi.Hello", "tjl", &bo)
	if err != nil {
		panic(err)
	}
	select {}
}

func RestApi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Debug on."})
}

type RpcApi struct {}

func (r *RpcApi) Hello(request string, response *string) error {
	*response = "hello: " + request
	return nil
}
