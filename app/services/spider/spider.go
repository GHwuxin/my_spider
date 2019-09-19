package spider

import (
	"fmt"
	"mypractice/spider/app/services/downloader"
	"mypractice/spider/app/services/pageprocesser"
	"mypractice/spider/app/services/pipeline"
	"mypractice/spider/app/services/resourcemanage"
	"mypractice/spider/app/services/scheduler"
	"mypractice/spider/models"
	"mypractice/spider/storage/logs"
	"time"
)

type Spider struct {
	Name           string
	Num            uint
	pageProcesser  pageprocesser.PageProcesser
	downloader     downloader.Downloader
	scheduler      scheduler.Scheduler
	piplelines     []pipeline.Pipeline
	resourceManage resourcemanage.ResourceManage
}

// Spider is scheduler module for all the other modules, like downloader, pipeline, scheduler and etc.
func NewSpider(name string, num uint, page pageprocesser.PageProcesser, down downloader.Downloader, sche scheduler.Scheduler, pip []pipeline.Pipeline) *Spider {

	spider := new(Spider)
	spider.Name = name
	if num <= 0 {
		spider.Num = 1
	} else {
		spider.Num = num
	}
	if page == nil {
		return nil
	}
	spider.pageProcesser = page
	if down == nil {
		spider.downloader = downloader.NewHttpDownloader()
	} else {
		spider.downloader = down
	}
	if sche == nil {
		spider.scheduler = scheduler.NewSimpleScheduler()
	} else {
		spider.scheduler = sche
	}
	if pip == nil {
		spider.piplelines = make([]pipeline.Pipeline, 0)
	} else {
		spider.piplelines = pip
	}

	return spider
}

func (this *Spider) Taskname() string {
	return this.Name
}

func (this *Spider) Threadnum() uint {
	return this.Num
}

func (this *Spider) RefreshTaskStatus(status *models.TaskStatus) {
	fmt.Println(fmt.Sprintf("当前task文件总个数：%d", status.FilesCount))
	fmt.Println(fmt.Sprintf("当前task文件下载索引：%d", status.CurrentIndex))
}

func (this *Spider) AddPipeline(p pipeline.Pipeline) *Spider {
	this.piplelines = append(this.piplelines, p)
	return this
}

func (this *Spider) AddUrl(url string, respType string) *Spider {
	req := models.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req)
	return this
}

func (this *Spider) AddUrlEx(url string, respType string, headerFile string, proxyHost string) *Spider {
	req := models.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	return this
}

func (this *Spider) AddUrlWithHeaderFile(url string, respType string, headerFile string) *Spider {
	req := models.NewRequestWithHeaderFile(url, respType, headerFile)
	this.AddRequest(req)
	return this
}

func (this *Spider) AddUrls(urls []string, respType string) *Spider {
	for _, url := range urls {
		req := models.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}
	return this
}

func (this *Spider) AddUrlsWithHeaderFile(urls []string, respType string, headerFile string) *Spider {
	for _, url := range urls {
		req := models.NewRequestWithHeaderFile(url, respType, headerFile)
		this.AddRequest(req)
	}
	return this
}

func (this *Spider) AddUrlsEx(urls []string, respType string, headerFile string, proxyHost string) *Spider {
	for _, url := range urls {
		req := models.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req.AddHeaderFile(headerFile).AddProxyHost(proxyHost))
	}
	return this
}

func (this *Spider) AddRequest(req *models.Request) *Spider {
	if req == nil {
		return this
	} else if req.Url() == "" {
		return this
	}
	this.scheduler.Put(req)
	return this
}

func (this *Spider) AddRequests(reqs []*models.Request) *Spider {
	for _, req := range reqs {
		this.AddRequest(req)
	}
	return this
}

func (this *Spider) One(url string, respType string) *models.PageItems {
	req := models.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	return this.OneByRequest(req)
}

func (this *Spider) All(urls []string, respType string) []*models.PageItems {
	for _, u := range urls {
		req := models.NewRequest(u, respType, "", "GET", "", nil, nil, nil, nil)
		this.AddRequest(req)
	}
	pip := pipeline.NewPipelineCollect()
	this.AddPipeline(pip)
	this.Run()
	return pip.GetCollected()
}

func (this *Spider) OneByRequest(req *models.Request) *models.PageItems {
	var reqs []*models.Request
	reqs = append(reqs, req)
	items := this.AllByRequestes(reqs)
	if len(items) != 0 {
		return items[0]
	}
	return nil
}

func (this *Spider) AllByRequestes(reqs []*models.Request) []*models.PageItems {
	for _, req := range reqs {
		this.AddRequest(req)
	}
	pip := pipeline.NewPipelineCollect()
	this.AddPipeline(pip)
	this.Run()
	return pip.GetCollected()
}

func (this *Spider) Run() {

	this.resourceManage = resourcemanage.NewResourceManageChan(this.Num)
	for {
		req := this.scheduler.Pop()
		if this.resourceManage.Has() == 0 && req == nil {
			this.pageProcesser.Finish()
			break
		} else if req == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		this.resourceManage.GetOne()
		go func(req *models.Request) {
			defer this.resourceManage.FreeOne()
			this.pageProcess(req)
		}(req)
	}
	this.close()
}

func (this *Spider) pageProcess(req *models.Request) {
	var p *models.Page
	defer func() {
		if err := recover(); err != nil {
			if errStr, ok := err.(string); ok {
				logs.Instance.Errorf("pageProcess panic error:%s", errStr)
			} else if err, ok := err.(error); ok {
				logs.Instance.Errorf("pageProcess panic error:%s", err.Error())
			}
		}
	}()
	// download page
	p = this.downloader.Download(req)
	if !p.Success() {
		logs.Instance.Errorf("get page error:%s", p.ErrorMessage())
		return
	}
	err := this.pageProcesser.Process(p)
	if err != nil {
		logs.Instance.Errorf("pageProcess error:%s", err.Error())
	}
	for _, req := range p.TargetRequests() {
		this.AddRequest(req)
	}
	// output
	if !p.Skip() {
		for _, pip := range this.piplelines {
			pip.Process(p.PageItems(), this)
		}
	}
}

func (this *Spider) close() {
	this.scheduler = scheduler.NewSimpleScheduler()
	this.downloader = downloader.NewHttpDownloader()
	this.piplelines = make([]pipeline.Pipeline, 0)
}
