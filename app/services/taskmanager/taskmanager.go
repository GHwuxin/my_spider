package taskmanager

import (
	"mypractice/spider/app/services/downloader"
	"mypractice/spider/app/services/pageprocesser"
	"mypractice/spider/app/services/pipeline"
	"mypractice/spider/app/services/resourcemanage"
	"mypractice/spider/app/services/scheduler"
	"mypractice/spider/app/services/spider"
	"mypractice/spider/config"
	"mypractice/spider/models"
	"time"
)

type TaskManager struct {
	Tasks []*config.TaskConfig
	rm    resourcemanage.ResourceManage
}

func NewTaskManager(num uint) *TaskManager {
	tm := new(TaskManager)
	tm.Tasks = make([]*config.TaskConfig, 0)
	tm.rm = resourcemanage.NewResourceManageChan(num)
	return tm
}

func (this *TaskManager) TaskNum() uint {
	return this.rm.Has()
}

func (this *TaskManager) AddTask(cfg ...*config.TaskConfig) *TaskManager {
	this.Tasks = append(this.Tasks, cfg...)
	return this
}

func (this *TaskManager) Run() {
	for _, task := range this.Tasks {
		this.rm.GetOne()
		go func(taskcfg *config.TaskConfig) {
			defer this.rm.FreeOne()
			page := pageprocesser.NewUrlPageProcesser(taskcfg.RequestCfg.Selecter.Name, taskcfg.RequestCfg.Selecter.Attr, taskcfg.RequestCfg.Selecter.Pattern)
			down := downloader.NewHttpDownloader()
			if len(taskcfg.RequestCfg.ChromedpTasks) > 0 {
				for _, arr := range taskcfg.RequestCfg.ChromedpTasks {
					down.AddChromedpTasks(arr)
				}
			}
			sche := scheduler.NewSimpleScheduler()
			pips := make([]pipeline.Pipeline, 0)
			pips = append(pips, pipeline.NewPipelineSave(taskcfg.DataSavePath))
			req := models.NewRequest(taskcfg.RequestCfg.Url, taskcfg.RequestCfg.ResponseType, "", taskcfg.RequestCfg.Method, taskcfg.RequestCfg.PostData, taskcfg.RequestCfg.Header, nil, nil, nil)
			if taskcfg.LoopTimespan > 0 {
				ticker := time.NewTicker(time.Second * time.Duration(taskcfg.LoopTimespan))
				for {
					select {
					case <-ticker.C:
						sp := spider.NewSpider(taskcfg.Name, taskcfg.ThreadNum, page, down, sche, pips)
						sp.AddRequest(req)
						sp.Run()
					}
				}
			} else {
				sp := spider.NewSpider(taskcfg.Name, taskcfg.ThreadNum, page, down, sche, pips)
				sp.AddRequest(req)
				sp.Run()
			}
		}(task)
	}
}
