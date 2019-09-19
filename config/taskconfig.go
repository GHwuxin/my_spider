package config

import (
	"mypractice/spider/utils"
	"os"
	"strings"
)

type TaskConfig struct {
	Name         string
	ThreadNum    uint
	DataSavePath string
	LoopTimespan uint
	RequestCfg   *RequestConfig
}

func NewTaskConfig() *TaskConfig {
	taskcfg := new(TaskConfig)
	taskcfg.Name = "default"
	taskcfg.ThreadNum = 1
	taskcfg.DataSavePath = "./download"
	taskcfg.RequestCfg = NewRequesetConfig()
	return taskcfg
}

func (this *TaskConfig) TestConfig() error {
	if this.Name == "" {
		this.Name = "default"
	}
	if this.ThreadNum <= 0 {
		this.ThreadNum = 1
	}
	if this.DataSavePath == "" {
		this.DataSavePath = "./download"
	}
	if !utils.IsExists(this.DataSavePath) &&
		!strings.Contains(this.DataSavePath, "yyyy") &&
		!strings.Contains(this.DataSavePath, "MM") &&
		!strings.Contains(this.DataSavePath, "dd") &&
		!strings.Contains(this.DataSavePath, "HH") &&
		!strings.Contains(this.DataSavePath, "mm") &&
		!strings.Contains(this.DataSavePath, "ss") {
		err := os.MkdirAll(this.DataSavePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := this.RequestCfg.TestConfig()
	return err
}
