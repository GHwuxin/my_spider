package models

import "time"

type DownloadFileStatus struct {
	FileUrl       string
	FilePath      string
	FileName      string
	FileDate      time.Time
	FileSize      int64
	FileWriteSize int64
	NetworkSpeed  int64
}
