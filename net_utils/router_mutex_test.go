package net_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/rpc"
	"testing"
)

func TestNewRouterMulti(t *testing.T) {
	multi := NewRouterMulti(9999)

	api := &Api{}
	multi.AddRpcHandle(api)
	multi.AddRestHandle(REST_GET, "/rest", api.Rest)
	multi.AddRestHandle(REST_GET, "/error", api.Error)

	MiddleList = BuildDefaultMiddleList()

	multi.Init()
	go multi.Loop()

	// 连接远程rpc
	rp, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	var bo *string
	err = rp.Call("Api.Hello", "tjl", &bo)
	if err != nil {
		panic(err)
	}
	select {}
}

type Api struct{}

func (r *Api) Hello(request string, response *string) error {
	*response = "hello: " + request
	return nil
}

func (r *Api) Rest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "This is rest api."})
}

func (r *Api) Error(c *gin.Context) {
	panic("My error")
}
