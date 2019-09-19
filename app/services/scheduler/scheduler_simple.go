package scheduler

import "mypractice/spider/models"

type SimpleScheduler struct {
	queue chan *models.Request
}

func NewSimpleScheduler() *SimpleScheduler {
	ch := make(chan *models.Request, 1024)
	return &SimpleScheduler{ch}
}

func (this *SimpleScheduler) Put(requ *models.Request) {
	this.queue <- requ
}

func (this *SimpleScheduler) Pop() *models.Request {
	if this.Count() > 0 {
		return <-this.queue
	} else {
		return nil
	}
}

func (this *SimpleScheduler) Count() int {
	return len(this.queue)
}
