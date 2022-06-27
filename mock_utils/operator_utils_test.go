package mock_utils

import (
	"fmt"
	"testing"
)

func RemoveHello(id, arg string) {
	fmt.Println("Remove: " + id + " " + arg)
}

func CreateHello(id, arg string) {
	fmt.Println("Create: " + id + " " + arg)
}

func DetailHello(id, arg string) {
	fmt.Println("Detail: " + id + " " + arg)
}

func SearchHello(id, arg string) {
	fmt.Println("Search: " + id + " " + arg)
}

func TestNewMockOPClient(t *testing.T) {

	funcArr := []func(string, string){RemoveHello, CreateHello, DetailHello, SearchHello,}

	generatCommandMap := GeneratCommandMapper(funcArr)

	// The SearchHello is loop func works.
	generatCommandMap["SearchHello"].CommandType = 1

	// Specifiy the RmeoveHello to be the flag to stop all group tasks.
	mockClient := NewMockOPClient("", "RemoveHello", generatCommandMap)

	mockClient.Loop("./")
}
