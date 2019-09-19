package models

type TaskCurrentStatus int

const (
	_ TaskCurrentStatus = iota - 3
	Init
	Fail
	Success
	Running
)

type TaskStatus struct {
	CurrentStatus   TaskCurrentStatus
	FilesCount      int
	SuccessFilesUrl []string
	FailFilesUrl    []string
	CurrentIndex    int
}

func NewTaskStatus() *TaskStatus {
	status := new(TaskStatus)
	status.CurrentStatus = Init
	status.FilesCount = 0
	status.SuccessFilesUrl = make([]string, 0)
	status.FailFilesUrl = make([]string, 0)
	status.CurrentIndex = 0
	return status
}
