package pageprocesser

import "mypractice/spider/models"

type PageProcesser interface {
	Process(p *models.Page) error
	Finish()
}
