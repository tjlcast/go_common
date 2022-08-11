package configer_utils

import (
	"encoding/json"
	"os"
	"path"
	"testing"
)

func TestParse(t *testing.T) {
	dir, _ := os.Getwd()
	t.Log(dir)

	configFilePath := path.Join(dir, "config", "app.conf")
	t.Log(configFilePath)

	// init
	_, _ = Parse(configFilePath)
	WatchConfig(configFilePath)

	s := GetGlobalV("runmode")
	if s != "dev" {
		t.Error("runmode is not dev")
	}

	v := GetSecV("dev", "name")
	if v != "hello" {
		t.Error("dev->name is not hello")
	}

	// get sec dev
	cs := GetSec("dev")
	if cs["name"] != "hello" {
		t.Error("Fail GetSec!")
	}

	// check
	bytes1, _ := json.Marshal(cs)
	t.Log(string(bytes1))
}
