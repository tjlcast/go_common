package net_utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tjlcast/go_common/log_utils"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	REST_GET    = "GET"
	REST_PUT    = "PUT"
	REST_POST   = "POST"
	REST_DELETE = "DELETE"

	/**
	After rpc call, Before others all.

	Default is nil
	If empty, do nothing;
	*/
	MiddleList []func(c *gin.Context)
)

/**
life cycle:
+ NewRouterMulti
+ AddXXX
+ Init
+ Loop
*/

type RouterMulti struct {
	port     int
	listener *net.Listener

	rpcHandles []interface{}
	// <method, <url, func>>
	restHandles map[string]map[string]gin.HandlerFunc

	RpcServer *rpc.Server
	Engine    *gin.Engine
}

func NewRouterMulti(port int) *RouterMulti {
	restMap := make(map[string]map[string]gin.HandlerFunc)
	restMap[REST_GET] = make(map[string]gin.HandlerFunc)
	restMap[REST_POST] = make(map[string]gin.HandlerFunc)
	restMap[REST_PUT] = make(map[string]gin.HandlerFunc)
	restMap[REST_DELETE] = make(map[string]gin.HandlerFunc)

	return &RouterMulti{
		port: port,
		//rpcHandles: []interface{}{},
		restHandles: restMap,
	}
}

func (rm *RouterMulti) IsEnable() bool {
	if rm.listener != nil {
		return true
	}
	return false
}

func (rm *RouterMulti) AddRpcHandle(handle interface{}) {
	rm.rpcHandles = append(rm.rpcHandles, handle)
}

func (rm *RouterMulti) AddRestHandle(method string, relativePath string, handler gin.HandlerFunc) {
	funcs, ok := rm.restHandles[method]
	if !ok {
		panic("Not support method: " + method)
	}
	_, ok = funcs[relativePath]
	if ok {
		panic("There has been defined url: " + relativePath)
	}
	funcs[relativePath] = handler
}

func (rm *RouterMulti) Init() []string {
	// listener
	endpoint := ":" + strconv.Itoa(rm.port)
	lis, err := GetOrBuildListener(endpoint)
	if err != nil {
		panic(err)
	}
	rm.listener = &lis

	// rpc
	infoCollect := []string{}
	rm.RpcServer = rpc.NewServer()
	for _, rpcHandle := range rm.rpcHandles {
		clz := reflect.TypeOf(rpcHandle)
		cname := clz.Elem().Name()
		for i := 0; i < clz.NumMethod(); i++ {
			methodName := clz.Method(i).Name
			sprintf := fmt.Sprintf("Rpc api: %s.%s\n", cname, methodName)
			infoCollect = append(infoCollect, sprintf)
		}
		err := rm.RpcServer.Register(rpcHandle)
		if err != nil {
			panic(err)
		}
	}
	infoCollect = append(infoCollect)

	// rest
	rm.Engine = gin.Default()
	rm.Engine.Use(func(ctx *gin.Context) {
		requestUrl := ctx.Request.RequestURI
		if strings.HasPrefix(requestUrl, "/_goRPC_") || ctx.Request.ProtoMajor == 2 &&
			strings.HasPrefix(ctx.GetHeader("Content-Type"), "application/grpc") {
			rm.RpcServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}
		ctx.Next()
	})
	// gin middle
	if MiddleList != nil && len(MiddleList) > 0 {
		for _, mid := range MiddleList {
			rm.Engine.Use(mid)
		}
	}
	funcs := rm.restHandles[REST_GET]
	for url, f := range funcs {
		rm.Engine.GET(url, f)

		sprintf := fmt.Sprintf("Rest api: [%s] %s\n", "GET", url)
		infoCollect = append(infoCollect, sprintf)
	}
	funcs = rm.restHandles[REST_POST]
	for url, f := range funcs {
		rm.Engine.POST(url, f)

		sprintf := fmt.Sprintf("Rest api: [%s] %s\n", "POST", url)
		infoCollect = append(infoCollect, sprintf)
	}
	funcs = rm.restHandles[REST_PUT]
	for url, f := range funcs {
		rm.Engine.PUT(url, f)

		sprintf := fmt.Sprintf("Rest api: [%s] %s\n", "PUT", url)
		infoCollect = append(infoCollect, sprintf)
	}
	funcs = rm.restHandles[REST_DELETE]
	for url, f := range funcs {
		rm.Engine.DELETE(url, f)

		sprintf := fmt.Sprintf("Rest api: [%s] %s\n", "DELETE", url)
		infoCollect = append(infoCollect, sprintf)
	}
	return infoCollect
}

func (rm *RouterMulti) Loop() {
	// rest
	err := rm.Engine.RunListener(*rm.listener)
	if err != nil {
		panic(err)
	}
	// rpc
	rm.RpcServer.Accept(*rm.listener)
}

func BuildDefaultMiddleList() []func(c *gin.Context) {
	return []func(c *gin.Context){Recover}
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log_utils.TestLog(fmt.Sprintf("panic: %v\n", r))

			stackInfo := string(debug.Stack())
			fmt.Println(stackInfo)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  errorToString(r),
				"data": stackInfo,
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
