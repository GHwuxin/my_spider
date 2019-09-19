package downloader

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mypractice/spider/models"
	"mypractice/spider/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html/charset"
)

type HttpDownloader struct {
	chromedpTasks [][]string
}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{chromedpTasks: make([][]string, 0)}
}

func (this *HttpDownloader) AddChromedpTasks(action ...[]string) *HttpDownloader {
	this.chromedpTasks = append(this.chromedpTasks, action...)
	return this
}

func (this *HttpDownloader) Download(req *models.Request) *models.Page {

	var p = models.NewPage(req)
	switch req.ResponceType() {
	case "html":
		return this.downloadHtml(p, req)
	case "json":
		fallthrough
	case "jsonp":
		return this.downloadJson(p, req)
	case "text":
		return this.downloadText(p, req)
	}
	return p
}

func (this *HttpDownloader) changeChromeHeadlessReadHtml(url string) (string, error) {

	if len(this.chromedpTasks) == 0 {
		return "", errors.New("chromedpTasks len is 0")
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()
	htmlStr := ""
	tasks := make(chromedp.Tasks, 0)
	for _, strArr := range this.chromedpTasks {
		if len(strArr) != 2 {
			continue
		}
		switch strArr[0] {
		case "Navigate":
			tasks = append(tasks, chromedp.Navigate(url))
		case "WaitReady":
			tasks = append(tasks, chromedp.WaitReady(strArr[1], chromedp.ByID))
		case "OuterHTML":
			tasks = append(tasks, chromedp.OuterHTML(strArr[1], &htmlStr, chromedp.ByQuery))
		}
	}
	err := chromedp.Run(ctx, tasks)
	return htmlStr, err
}

// Charset auto determine. Use golang.org/x/net/html/charset. Get page body and change it to utf-8
func (this *HttpDownloader) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) (string, error) {

	destReader, err := charset.NewReader(sor, contentTypeStr)
	if err != nil {
		destReader = sor
	}
	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		return "", err
	}
	bodystr := string(sorbody)
	return bodystr, nil
}

func (this *HttpDownloader) changeCharsetEncodingAutoGzipSupport(contentTypeStr string, sor io.ReadCloser) (string, error) {

	gzipReader, err := gzip.NewReader(sor)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()
	destReader, err := charset.NewReader(gzipReader, contentTypeStr)
	if err != nil {
		destReader = sor
	}
	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		return "", err
	}
	bodystr := string(sorbody)
	return bodystr, nil
}

// Download file and change the charset of page charset.
func (this *HttpDownloader) downloadFile(p *models.Page, req *models.Request) (*models.Page, string) {
	var err error
	var urlstr string
	if urlstr = req.Url(); len(urlstr) == 0 {
		p.SetStatus(true, "url is empty")
		return p, ""
	}
	var resp *http.Response
	if proxystr := req.ProxyHost(); len(proxystr) != 0 {
		resp, err = connectByHttpProxy(p, req)
	} else {
		resp, err = connectByHttp(p, req)
	}
	if err != nil {
		p.SetStatus(true, err.Error())
		return p, ""
	}
	p.Header = resp.Header
	p.Cookies = resp.Cookies()

	bodyStr := ""
	if len(this.chromedpTasks) > 0 {
		bodyStr, err = this.changeChromeHeadlessReadHtml(urlstr)
	} else if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyStr, err = this.changeCharsetEncodingAutoGzipSupport(resp.Header.Get("Content-Type"), resp.Body)
	} else {
		bodyStr, err = this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
	}
	defer resp.Body.Close()
	if err != nil {
		p.SetStatus(true, err.Error())
	}
	return p, bodyStr
}

func (this *HttpDownloader) downloadHtml(p *models.Page, req *models.Request) *models.Page {
	var err error
	p, destbody := this.downloadFile(p, req)
	if !p.Success() {
		return p
	}
	bodyReader := bytes.NewReader([]byte(destbody))
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
		p.SetStatus(true, err.Error())
		return p
	}
	var body string
	if body, err = doc.Html(); err != nil {
		p.SetStatus(true, err.Error())
		return p
	}
	p.SetBody(body).SetHtmlParser(doc).SetStatus(false, "")
	return p
}

func (this *HttpDownloader) downloadJson(p *models.Page, req *models.Request) *models.Page {
	var err error
	p, destbody := this.downloadFile(p, req)
	if !p.Success() {
		return p
	}
	var body []byte
	body = []byte(destbody)
	mtype := req.ResponceType()
	if mtype == "jsonp" {
		tmpstr := utils.JsonpToJson(destbody)
		body = []byte(tmpstr)
	}
	var r *simplejson.Json
	if r, err = simplejson.NewJson(body); err != nil {
		p.SetStatus(true, err.Error())
		return p
	}
	// json result
	p.SetBody(string(body)).SetJson(r).SetStatus(false, "")
	return p
}

func (this *HttpDownloader) downloadText(p *models.Page, req *models.Request) *models.Page {
	p, destbody := this.downloadFile(p, req)
	if !p.Success() {
		return p
	}
	p.SetBody(destbody).SetStatus(false, "")
	return p
}

// choose http GET/method to download
func connectByHttp(p *models.Page, req *models.Request) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: req.RedirectFunc(),
	}
	httpreq, err := http.NewRequest(req.Method(), req.Url(), strings.NewReader(req.Postdata()))
	if err != nil {
		return nil, err
	}
	if header := req.Header(); header != nil {
		httpreq.Header = req.Header()
	}
	if cookies := req.Cookies(); cookies != nil {
		for i := range cookies {
			httpreq.AddCookie(cookies[i])
		}
	}
	var resp *http.Response
	if resp, err = client.Do(httpreq); err != nil {
		if e, ok := err.(*url.Error); ok && e.Err != nil && e.Err.Error() == "normal" {
			//  normal
		} else {
			return nil, err
		}
	}
	return resp, nil
}

// choose a proxy server to excute http GET/method to download
func connectByHttpProxy(p *models.Page, in_req *models.Request) (*http.Response, error) {
	request, _ := http.NewRequest("GET", in_req.Url(), nil)
	proxy, err := url.Parse(in_req.ProxyHost())
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
