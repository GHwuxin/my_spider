package pipeline

import (
	"mypractice/spider/app/services/task"
	"mypractice/spider/models"
)

type PipelineCollect struct {
	collector []*models.PageItems
}

func NewPipelineCollect() *PipelineCollect {
	collector := make([]*models.PageItems, 0)
	return &PipelineCollect{collector: collector}
}

func (this *PipelineCollect) Process(items *models.PageItems, t task.Task) {
	this.collector = append(this.collector, items)
}

func (this *PipelineCollect) GetCollected() []*models.PageItems {
	return this.collector
}
