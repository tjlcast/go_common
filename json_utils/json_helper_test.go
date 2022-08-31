package json_utils

import (
	"fmt"
	"testing"
)

var (
	jsonObjStr = `{
"name": "tjlcast",
"age": 123
}`

	jsonArrStr = `[{"Name": "tjl", "age": 24}, {"Name": "wj", "age":21}]`
)

func TestJsonStrTransfer(t *testing.T) {

}

func TestParseJsonArr(t *testing.T) {
	arr := ParseJsonArr(jsonArrStr)
	for _, a := range arr {
		ele := a.(map[string]interface{})
		fmt.Println(ele["Name"], " - ", ele["age"])
	}
}

func TestParseJsonObj(t *testing.T) {
	obj := ParseJsonObj(jsonObjStr)
	for key, val := range obj {
		fmt.Println(key, " ", val)
	}
}

type Student struct {
	Name string
	Score []*Class
}

type Class struct {
	Name string
	Score int
}

func TestToJsonStr(t *testing.T) {
	student := Student{"tjl", []*Class{&Class{"math", 98}, &Class{"Chinese", 67}}}
	str := ToJsonStr(student)
	fmt.Println(str)
}
