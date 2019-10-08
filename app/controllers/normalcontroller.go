package controllers

import (
	"errors"
	"fmt"
	"mypractice/spider/app/services/filesave"
	"mypractice/spider/app/services/resourcemanage"
	"mypractice/spider/models"
	"mypractice/spider/storage/logs"
	"mypractice/spider/utils"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type NormalController struct {
	RootDir        string
	STime          string
	ETime          string
	resourceManage resourcemanage.ResourceManage
}

func NewNormalController(rootDir, sTime, eTime string) *NormalController {

	normalC := new(NormalController)
	normalC.ETime = eTime
	normalC.STime = sTime
	normalC.RootDir = rootDir
	return normalC
}

func (normalC *NormalController) Run() error {

	sDatetime, err := time.Parse("2006-01-02 15:04:05", normalC.STime)
	if err != nil {
		return errors.New("stime parse error:" + err.Error())
	}
	eDatetime, err := time.Parse("2006-01-02 15:04:05", normalC.ETime)
	if err != nil {
		return errors.New("etime parse error:" + err.Error())
	}
	sTimeZeroStr := sDatetime.Format("2006-01-02") + " 00:00:00"
	eTimeZeroStr := eDatetime.Add(time.Hour*24).Format("2006-01-02") + " 00:00:00"
	sTimeZero, _ := time.Parse("2006-01-02 15:04:05", sTimeZeroStr)
	eTimeZero, _ := time.Parse("2006-01-02 15:04:05", eTimeZeroStr)
	tempTime := sTimeZero
	normalC.resourceManage = resourcemanage.NewResourceManageChan(8)
	for tempTime.Before(eTimeZero) {
		normalC.resourceManage.GetOne()
		if !tempTime.After(eDatetime) && !tempTime.Before(sDatetime) {
			go func(tempDir string, tempT time.Time) {
				defer normalC.resourceManage.FreeOne()
				// http://d1.weather.com.cn/radar_channel/radar/pic/ACHN_QREF_20191008_064000.png
				filename := fmt.Sprintf("ACHN_QREF_%s_%s.png", tempT.Format("20060102"), tempT.Format("150405"))
				url := fmt.Sprintf("http://d1.weather.com.cn/radar_channel/radar/pic/%s", filename)
				tempDir = strings.Replace(tempDir, "yyyy", tempT.Format("2006"), -1)
				tempDir = strings.Replace(tempDir, "MM", tempT.Format("01"), -1)
				tempDir = strings.Replace(tempDir, "dd", tempT.Format("02"), -1)
				tempDir = strings.Replace(tempDir, "HH", tempT.Format("15"), -1)
				tempDir = strings.Replace(tempDir, "mm", tempT.Format("04"), -1)
				tempDir = strings.Replace(tempDir, "ss", tempT.Format("05"), -1)
				var err error
				if !utils.IsExists(tempDir) {
					err = os.MkdirAll(tempDir, os.ModePerm)
					if err != nil {
						logs.Instance.Warnf("%s make dir error:%s", url, err.Error())
					}
				}
				if err == nil {
					tempFilepath := path.Join(tempDir, filename)
					header := make(http.Header)
					header.Add("Referer", "http://www.weather.com.cn/radar/radar.html")
					err = filesave.HttpDownloadFile(url, header, nil, tempFilepath, normalC.RefreshDownloadStatus)
					if err != nil {
						logs.Instance.Warnf("download file(%s) fail error:%s", tempFilepath, err.Error())
					}
				}
			}(normalC.RootDir, tempTime)

		}
		tempTime = tempTime.Add(time.Minute * 10)
	}
	return nil
}

func (normalC *NormalController) RefreshDownloadStatus(status *models.DownloadFileStatus) {
	fmt.Println(fmt.Sprintf("当前文件URL：%s", status.FileUrl))
}
