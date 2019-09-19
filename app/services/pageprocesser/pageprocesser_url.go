package pageprocesser

import (
	"errors"
	"mypractice/spider/models"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type UrlPageProcesser struct {
	selecter string
	attr     string
	pattern  string
}

func NewUrlPageProcesser(selecter, attr, pattern string) *UrlPageProcesser {
	return &UrlPageProcesser{selecter: selecter, attr: attr, pattern: pattern}
}

// <img src="http://d1.weather.com.cn/radar_channel/radar/pic/ACHN_QREF_20190909_155000.png" style="width: 100%;">
func (this *UrlPageProcesser) Process(p *models.Page) error {
	if !p.Success() {
		return errors.New(p.ErrorMessage())
	}
	query := p.HtmlParser()
	query.Find(this.selecter).Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr(this.attr)
		regexpPng := regexp.MustCompile(this.pattern)
		if regexpPng.MatchString(src) {
			arrStr := strings.Split(src, "/")
			tempFilename := arrStr[len(arrStr)-1]
			p.AddField(tempFilename, src)
		}
	})
	return nil
}

func (this *UrlPageProcesser) Finish() {

}
