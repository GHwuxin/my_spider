package config

import "net/http"

type RequestConfig struct {
	Url           string
	ResponseType  string      // Responce type: html json jsonp text
	Method        string      // GET POST
	PostData      string      // POST data
	Header        http.Header // http header
	Selecter      *SelecterConfig
	ChromedpTasks [][]string // chromedp actions
}

func NewRequesetConfig() *RequestConfig {

	req := new(RequestConfig)
	req.ResponseType = "html"
	req.Method = "GET"
	req.Selecter = NewSelecterConfig()
	req.Header = make(http.Header)
	req.ChromedpTasks = make([][]string, 0)
	return req
}

func (this *RequestConfig) TestConfig() error {
	if this.ResponseType == "" {
		this.ResponseType = "html"
	}
	if this.Method == "" {
		this.Method = "GET"
	}
	if this.Selecter == nil {
		this.Selecter = NewSelecterConfig()
	}
	if this.Header == nil {
		this.Header = make(http.Header)
	}
	if this.ChromedpTasks == nil {
		this.ChromedpTasks = make([][]string, 0)
	}
	return nil
}
