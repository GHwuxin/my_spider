package models

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
)

type Page struct {
	isfail         bool              // The isfail is true when crawl process is failed and errormsg is the fail resean.
	errormsg       string            // error message
	Request        *Request          // The request is crawled by spider that contains url and relevent information.
	body           string            // The body is plain text of crawl result.
	Header         http.Header       // header
	Cookies        []*http.Cookie    // cookie
	docParser      *goquery.Document // The docParser is a pointer of goquery boject that contains html result.
	jsonMap        *simplejson.Json  // The jsonMap is the json result.
	pItems         *PageItems        // The pItems is object for save Key-Values in PageProcesser.And pItems is output in Pipline.
	targetRequests []*Request        // The targetRequests is requests to put into Scheduler.
}

// NewPage returns initialized Page object.
func NewPage(req *Request) *Page {
	return &Page{pItems: NewPageItems(req), Request: req}
}

// SetStatus save status info about download process.
func (this *Page) SetStatus(isfail bool, errormsg string) {
	this.isfail = isfail
	this.errormsg = errormsg
}

// SetSkip set label "skip" of PageItems.
// PageItems will not be saved in Pipeline wher skip is set true
func (this *Page) SetSkip(skip bool) {
	this.pItems.SetSkip(skip)
}

// SetBody saves plain string crawled in Page.
func (this *Page) SetBody(body string) *Page {
	this.body = body
	return this
}

// SetHtmlParser saves goquery object binded to target crawl result.
func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
	this.docParser = doc
	return this
}

// HtmlParser returns goquery object binded to target crawl result.
func (this *Page) ResetHtmlParser() *goquery.Document {
	r := strings.NewReader(this.body)
	var err error
	this.docParser, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil
	}
	return this.docParser
}

// SetJson saves json result.
func (this *Page) SetJson(js *simplejson.Json) *Page {
	this.jsonMap = js
	return this
}

// Success test whether download process success or not.
func (this *Page) Success() bool {
	return !this.isfail
}

// ErrorMessage show the download error message.
func (this *Page) ErrorMessage() string {
	return this.errormsg
}

// PageItems returns PageItems object that record KV pair parsed in PageProcesser.
func (this *Page) PageItems() *PageItems {
	return this.pItems
}

// Skip returns skip label of PageItems.
func (this *Page) Skip() bool {
	return this.pItems.Skip()
}

// UrlTag returns name of url.
func (this *Page) UrlTag() string {
	return this.Request.UrlTag()
}

// TargetRequests returns the target requests that will put into Scheduler
func (this *Page) TargetRequests() []*Request {
	return this.targetRequests
}

// HtmlParser returns goquery object binded to target crawl result.
func (this *Page) HtmlParser() *goquery.Document {
	return this.docParser
}

// Json returns json result.
func (this *Page) Json() *simplejson.Json {
	return this.jsonMap
}

// AddField saves KV string pair to PageItems preparing for Pipeline
func (this *Page) AddField(key string, value string) {
	this.pItems.AddItem(key, value)
}

// AddTargetRequest adds one new Request waitting for crawl.
func (this *Page) AddTargetRequest(url string, respType string) *Page {
	this.targetRequests = append(this.targetRequests, NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil))
	return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
	for _, url := range urls {
		this.AddTargetRequest(url, respType)
	}
	return this
}

// AddTargetRequestWithProxy adds one new Request waitting for crawl.
func (this *Page) AddTargetRequestWithProxy(url string, respType string, proxyHost string) *Page {
	this.targetRequests = append(this.targetRequests, NewRequestWithProxy(url, respType, "", "GET", "", nil, nil, proxyHost, nil, nil))
	return this
}

// AddTargetRequestsWithProxy adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithProxy(urls []string, respType string, proxyHost string) *Page {
	for _, url := range urls {
		this.AddTargetRequestWithProxy(url, respType, proxyHost)
	}
	return this
}

// AddTargetRequest adds one new Request with header file for waitting for crawl.
func (this *Page) AddTargetRequestWithHeaderFile(url string, respType string, headerFile string) *Page {
	this.targetRequests = append(this.targetRequests, NewRequestWithHeaderFile(url, respType, headerFile))
	return this
}

// AddTargetRequest adds one new Request waitting for crawl.
// The respType is "html" or "json" or "jsonp" or "text".
// The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
// The method is POST or GET.
// The postdata is http body string.
// The header is http header.
// The cookies is http cookies.
func (this *Page) AddTargetRequestWithParams(req *Request) *Page {
	this.targetRequests = append(this.targetRequests, req)
	return this
}

// AddTargetRequests adds new Requests waitting for crawl.
func (this *Page) AddTargetRequestsWithParams(reqs []*Request) *Page {
	for _, req := range reqs {
		this.AddTargetRequestWithParams(req)
	}
	return this
}
