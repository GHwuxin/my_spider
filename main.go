package main

import (
	"encoding/json"
	"io/ioutil"
	"mypractice/spider/app/services/taskmanager"
	"mypractice/spider/config"
	"mypractice/spider/storage/logs"
	"mypractice/spider/utils"
	"net/http"
	"os"
	"path"
)

func main() {

	dir := "./config/tasks"
	if !utils.IsExists(dir) {
		return
	}
	files, err := utils.AllFiles(dir)
	if err != nil {
		logs.Instance.Warnf("get ./config/tasks all cfg error:%s", err.Error())
	}
	taskM := taskmanager.NewTaskManager(1)
	for _, file := range files {
		if !utils.IsExists(file) {
			logs.Instance.Warnf("not find file:%s", file)
			continue
		}
		if path.Ext(file) != ".josn" {
			continue
		}
		cfgEntry := config.NewTaskConfig()
		cfgFile, err := os.Open(file)
		if err != nil {
			logs.Instance.Warnf("open file(%s) error:%s", file, err.Error())
			continue
		}
		cfgBytes, err := ioutil.ReadAll(cfgFile)
		err = json.Unmarshal(cfgBytes, &cfgEntry)
		if err != nil {
			logs.Instance.Warnf("Unmarshal file(%s) error:%s", file, err.Error())
			continue
		}
		err = cfgEntry.TestConfig()
		if err != nil {
			logs.Instance.Warnf("test file(%s) error:%s", file, err.Error())
			continue
		}
		taskM.AddTask(cfgEntry)
	}
	taskM.Run()

	mux := http.NewServeMux()
	err = http.ListenAndServe(":6666", mux)
}
