package filesave

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mypractice/spider/models"
	"mypractice/spider/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

func HttpDownloadFile(url string, header http.Header, cookies []*http.Cookie, localFilePath string, refreshStatus func(*models.DownloadFileStatus)) error {

	if utils.IsExists(localFilePath) {
		return errors.New("local file is exists!")
	}
	var buff = make([]byte, 32*1024)
	tempFilePath := localFilePath + ".download"
	// TODO:断点续传
	if utils.IsExists(tempFilePath) {
		err := os.Remove(tempFilePath)
		if err != nil {
			return errors.New("DownloadFile error:" + err.Error())
		}
	}

	defer func() {
		if utils.IsExists(tempFilePath) {
			os.Remove(tempFilePath)
		}
	}()
	tempFile, err := os.OpenFile(tempFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("open or create " + tempFilePath + " error:" + err.Error())
	}
	defer tempFile.Close()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if header != nil {
		for key, value := range header {
			for _, item := range value {
				req.Header.Add(key, item)
			}
		}
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Accept-Encoding", "identity") // 强迫服务器返回非压缩格式
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response StatusCode is %d", resp.StatusCode))
	}
	fileSize := resp.ContentLength
	if fileSize <= 0 {
		return errors.New(fmt.Sprintf("response ContentLength is %d", fileSize))
	}
	reader := bufio.NewReaderSize(resp.Body, 1024*32)
	writer := bufio.NewWriter(tempFile)
	var endRead bool
	var written int
	go func() error {
		for {
			lenR, errR := reader.Read(buff)
			if lenR > 0 {
				lenW, errW := writer.Write(buff[0:lenR])
				if lenW > 0 {
					written += lenW
				}
				if errW != nil {
					err = errW
					break
				}
				if lenR != lenW {
					err = io.ErrShortWrite
					break
				}
			}
			if errR != nil {
				if errR != io.EOF && lenR == 0 {
					err = errR
				} else {
					endRead = true
				}
				break
			}
		}
		if err != nil {
			return err
		}
		return nil
	}()

	arrTemp := strings.Split(localFilePath, "/")
	filename := arrTemp[len(arrTemp)-1]
	downloadStatus := new(models.DownloadFileStatus)
	downloadStatus.FileUrl = url
	downloadStatus.FilePath = localFilePath
	downloadStatus.FileSize = fileSize
	downloadStatus.FileName = filename
	downloadStatus.FileDate = utils.ParseDatetime(filename)
	spaceTime := time.Second * 1
	ticker := time.NewTicker(spaceTime)
	lastWtn := 0
	stop := false
	for {
		select {
		case <-ticker.C:
			speed := written - lastWtn
			if (written-lastWtn == 0) && endRead {
				ticker.Stop()
				stop = true
				writer.Flush()
				break
			} else {
				downloadStatus.FileWriteSize = int64(written)
				downloadStatus.NetworkSpeed = int64(speed)
			}
			lastWtn = written
		}
		if stop {
			break
		}
		refreshStatus(downloadStatus)
	}
	if err == nil {
		tempFile.Close()
		err = os.Rename(tempFilePath, localFilePath)
	}
	return err
}
