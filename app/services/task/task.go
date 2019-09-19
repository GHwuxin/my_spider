package task

import "mypractice/spider/models"

// The Task represents interface that contains environment variables.
// It inherits by Spider.
type Task interface {
	Taskname() string
	Threadnum() uint
	RefreshTaskStatus(status *models.TaskStatus)
}
