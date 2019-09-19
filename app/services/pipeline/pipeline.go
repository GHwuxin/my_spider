package pipeline

import (
	"mypractice/spider/app/services/task"
	"mypractice/spider/models"
)

type Pipeline interface {
	// The Process implements result persistent.
	// The items has the result be crawled.
	// The t has informations of this crawl task.
	Process(items *models.PageItems, t task.Task)
}

// The interface CollectPipeline recommend result in process's memory temporarily.
type CollectPipeline interface {
	Pipeline
	GetCollected() []*models.PageItems // The GetCollected returns result saved in in process's memory temporarily.
}
