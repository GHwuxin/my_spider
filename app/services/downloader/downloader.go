package downloader

import "mypractice/spider/models"

type Downloader interface {
	Download(req *models.Request) *models.Page
}
