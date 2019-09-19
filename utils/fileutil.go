package utils

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func AllFiles(dirPth string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			tempFiles, err := AllFiles(dirPth + PthSep + fi.Name())
			if err != nil {
				return nil, err
			}
			for _, f := range tempFiles {
				files = append(files, f)
			}
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func ParseDatetime(fileName string) time.Time {
	re := regexp.MustCompile("\\d{4,14}")
	dateStr := re.FindString(strings.Replace(strings.Replace(fileName, "_", "", -1), "-", "", -1))
	fileTime := time.Now()
	layout := "20060102150405"[:len(dateStr)]
	fileTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Now()
	}
	return fileTime
}
