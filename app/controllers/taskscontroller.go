package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mypractice/spider/app/services/taskmanager"
	"mypractice/spider/config"
	"mypractice/spider/storage/logs"
	"mypractice/spider/utils"
	"os"
	"path"
)

type TasksController struct {
	TasksCfgDir  string
	TasksSyncNum uint
}

func NewTasksController(cfgDir string, num uint) *TasksController {

	taskC := new(TasksController)
	taskC.TasksSyncNum = num
	taskC.TasksCfgDir = cfgDir
	return taskC
}

func (tasksC *TasksController) Run() error {

	if !utils.IsExists(tasksC.TasksCfgDir) {
		return errors.New("this dir is not exists!")
	}
	files, err := utils.AllFiles(tasksC.TasksCfgDir)
	if err != nil {
		logs.Instance.Warnf("get ./config/tasks all cfg error:%s", err.Error())
	}
	if tasksC.TasksSyncNum == 0 {
		tasksC.TasksSyncNum = 1
	}
	taskM := taskmanager.NewTaskManager(tasksC.TasksSyncNum)
	for _, file := range files {
		if !utils.IsExists(file) {
			logs.Instance.Warnf("not find file:%s", file)
			continue
		}
		if path.Ext(file) != ".json" {
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
	return nil
}
