package scheduler

import (
	"mypractice/spider/models"
)

type Scheduler interface {
	Put(requ *models.Request)
	Pop() *models.Request
	Count() int
}
