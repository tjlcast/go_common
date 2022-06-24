package configer_utils

import (
	"fmt"
	"log"

	"github.com/Unknwon/goconfig"
	"github.com/fsnotify/fsnotify"
)

const (
	DEFAULT_SECTION = goconfig.DEFAULT_SECTION
)

var cfgMap map[string]map[string]string = make(map[string]map[string]string)
var filePath string

func Parse(fpath string) (c map[string]map[string]string, err error) {
	cfg, err := goconfig.LoadConfigFile(fpath)
	if err != nil {
		return
	}
	filePath = fpath

	sections := cfg.GetSectionList()
	for _, v := range sections {
		cfgMap[v] = make(map[string]string, 0)
		keys := cfg.GetKeyList(v)
		for _, b := range keys {
			cfgMap[v][b], _ = cfg.GetValue(v, b)
		}
	}
	return cfgMap, err
}

func GetAllCfg() (c map[string]map[string]string) {
	return cfgMap
}

func GetSec(key string) (c map[string]string) {
	return cfgMap[key]
}

func ReloadAllCfg() (c map[string]map[string]string, err error) {
	return Parse(filePath)
}

func WatchConfig(filepath ...string) {

	go func() {
		watch, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		defer watch.Close()

		for _, fpath := range filepath {
			err = watch.Add(fpath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("WatchConfig " + fpath)
		}

		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						ReloadAllCfg()
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()
}


