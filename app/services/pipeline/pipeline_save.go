package pipeline

import (
	"fmt"
	"mypractice/spider/app/services/filesave"
	"mypractice/spider/app/services/resourcemanage"
	"mypractice/spider/app/services/task"
	"mypractice/spider/models"
	"mypractice/spider/storage/logs"
	"mypractice/spider/utils"
	"os"
	"path"
	"strings"
)

type DownloadResult struct {
	url string
	err error
}

type PipelineSave struct {
	rootDir        string
	resourceManage resourcemanage.ResourceManage
	results        chan *DownloadResult
}

func NewPipelineSave(path string) *PipelineSave {
	isExists := utils.IsExists(path)
	if !isExists &&
		!strings.Contains(path, "yyyy") &&
		!strings.Contains(path, "MM") &&
		!strings.Contains(path, "dd") &&
		!strings.Contains(path, "HH") &&
		!strings.Contains(path, "mm") &&
		!strings.Contains(path, "ss") {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			tempPath := "./tmp"
			os.MkdirAll(tempPath, os.ModePerm)
			return &PipelineSave{rootDir: tempPath}
		}
	}
	return &PipelineSave{rootDir: path}
}

func (this *PipelineSave) Process(items *models.PageItems, t task.Task) {

	taskStatus := models.NewTaskStatus()
	taskStatus.FilesCount = len(items.All())
	this.results = make(chan *DownloadResult, len(items.All()))
	this.resourceManage = resourcemanage.NewResourceManageChan(t.Threadnum())
	defer func() {
		if err := recover(); err != nil {
			if errStr, ok := err.(string); ok {
				logs.Instance.Errorf("pipeline panic error:%s", errStr)
			} else if err, ok := err.(error); ok {
				logs.Instance.Errorf("pipeline panic error:%s", err.Error())
			}
		}
	}()
	go func() {
		for key, value := range items.All() {
			this.resourceManage.GetOne()
			taskStatus.CurrentStatus = models.Running
			go func(tempDir, filename, url string, results chan *DownloadResult) {
				defer this.resourceManage.FreeOne()
				fileTime := utils.ParseDatetime(filename)
				tempDir = strings.Replace(tempDir, "yyyy", fileTime.Format("2006"), -1)
				tempDir = strings.Replace(tempDir, "MM", fileTime.Format("01"), -1)
				tempDir = strings.Replace(tempDir, "dd", fileTime.Format("02"), -1)
				tempDir = strings.Replace(tempDir, "HH", fileTime.Format("15"), -1)
				tempDir = strings.Replace(tempDir, "mm", fileTime.Format("04"), -1)
				tempDir = strings.Replace(tempDir, "ss", fileTime.Format("05"), -1)
				var err error
				if !utils.IsExists(tempDir) {
					err = os.MkdirAll(tempDir, os.ModePerm)
				}
				if err != nil {
					results <- &DownloadResult{url: url, err: err}
					logs.Instance.Errorf("create dir error:%s", err.Error())
				} else {
					tempFilepath := path.Join(tempDir, filename)
					header := items.Request().Header()
					isHave := false
					for _, v := range header["Referer"] {
						if v == items.Request().Url() {
							isHave = true
							break
						}
					}
					if !isHave {
						header.Add("Referer", items.Request().Url())
					}
					err = filesave.HttpDownloadFile(url, header, nil, tempFilepath, RefreshDownloadStatus)
					results <- &DownloadResult{url: url, err: err}
					if err != nil {
						logs.Instance.Warnf("download file(%s) fail error:", tempFilepath, err.Error())
					}
				}
			}(this.rootDir, key, value, this.results)
		}
	}()

	// wait for gorutiount exit
	for i := 0; i < len(items.All()); i++ {
		result := <-this.results
		if result.err != nil {
			taskStatus.FailFilesUrl = append(taskStatus.FailFilesUrl, result.url)
		} else {
			taskStatus.SuccessFilesUrl = append(taskStatus.SuccessFilesUrl, result.url)
		}
		taskStatus.CurrentIndex++
		if i == len(items.All())-1 {
			if len(taskStatus.SuccessFilesUrl) == 0 {
				taskStatus.CurrentStatus = models.Fail
			} else {
				taskStatus.CurrentStatus = models.Success
			}
		}
		t.RefreshTaskStatus(taskStatus)
	}
}

func RefreshDownloadStatus(status *models.DownloadFileStatus) {
	process := float64(status.FileWriteSize) / float64(status.FileSize) * 100
	if process > 100 {
		process = 100
	}
	fmt.Println(fmt.Sprintf("当前文件URL：%s", status.FileUrl))
	fmt.Println(fmt.Sprintf("当前文件本地路径：%s", status.FilePath))
	fmt.Println(fmt.Sprintf("当前文件名称：%s", status.FileName))
	fmt.Println(fmt.Sprintf("当前文件时间：%s", status.FileDate.Format("2006-01-02 15:04:05")))
	fmt.Println(fmt.Sprintf("当前文件大小：%s", utils.BytesToSize(status.FileSize)))
	fmt.Println(fmt.Sprintf("当前文件进度：%.2f", process) + "%")
	fmt.Println(fmt.Sprintf("当前下载网速：%s", utils.BytesToSize(status.NetworkSpeed)))
}
