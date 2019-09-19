package logs

import (
	"mypractice/spider/utils"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"

	"github.com/sirupsen/logrus"
)

var Instance *logrus.Logger
var Error error
var once sync.Once

func init() {
	once.Do(func() {
		logDir := "./storage/logs"
		if utils.IsExists(logDir) {
			err := os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				Error = err
			}
		}
		logPath := filepath.Join(logDir, "log")
		Instance = logrus.New()
		Instance.SetReportCaller(true)
		Instance.SetLevel(logrus.InfoLevel)
		Instance.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
		writer, err := rotatelogs.New(
			logPath+".%Y%m%d%H%M%S",
			rotatelogs.WithLinkName(logPath),
			rotatelogs.WithRotationTime(time.Hour),
			rotatelogs.WithMaxAge(time.Hour*24*3),
		)

		// lfsHook := lfshook.NewHook(lfshook.WriterMap{
		// 	logrus.DebugLevel: writer,
		// 	logrus.InfoLevel:  writer,
		// 	logrus.WarnLevel:  writer,
		// 	logrus.ErrorLevel: writer,
		// 	logrus.FatalLevel: writer,
		// 	logrus.PanicLevel: writer,
		// }, &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
		// Instance.AddHook(lfsHook)

		if err != nil {
			Error = err
		}
		Instance.SetOutput(writer)
	})
}
