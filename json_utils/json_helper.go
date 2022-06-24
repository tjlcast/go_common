package json_utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseJsonObj(jsonStr string) (map[string]interface{}) {
	var p map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err!=nil {
		fmt.Printf("Fail to parse map: %s.\n", jsonStr)
	}
	return p
}

func ParseJsonArr(jsonStr string) ([]interface{}) {
	var p []interface{}
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err!=nil {
		fmt.Printf("Fail to parse arr: %s.\n", jsonStr)
	}
	return p
}

func ToJsonStr(obj interface{}) string {
	bytes, e := json.Marshal(obj)
	if e != nil {
		fmt.Printf("obj is not a json")
	}
	return string(bytes)
}

func JsonStrTransfer(jsonStr string) string {
	transferred := strings.Trim(jsonStr, "\"")
	transferred = strings.ReplaceAll(transferred, "\\", "")
	return transferred
}