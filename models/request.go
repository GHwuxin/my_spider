package models

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

type Request struct {
	url           string
	respType      string         // Responce type: html json jsonp text
	method        string         // GET POST
	postData      string         // POST data
	urlTag        string         // name for marking url and distinguish different urls in PageProcesser and Pipeline
	header        http.Header    // http header
	cookies       []*http.Cookie // http cookies
	proxyHost     string         //proxy host   example='localhost:80'
	checkRedirect func(req *http.Request, via []*http.Request) error
	meta          interface{}
}

func NewRequest(
	url string,
	respType string,
	urltag string,
	method string,
	postdata string,
	header http.Header,
	cookies []*http.Cookie,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {

	return &Request{url, respType, method, postdata, urltag, header, cookies, "", checkRedirect, meta}
}

func NewRequestWithProxy(
	url string,
	respType string,
	urltag string,
	method string,
	postdata string,
	header http.Header,
	cookies []*http.Cookie,
	proxyHost string,
	checkRedirect func(req *http.Request, via []*http.Request) error,
	meta interface{}) *Request {

	return &Request{url, respType, method, postdata, urltag, header, cookies, proxyHost, checkRedirect, meta}
}

func NewRequestWithHeaderFile(url string, respType string, headerFile string) *Request {

	_, err := os.Stat(headerFile)
	if err != nil {
		//file is not exist , using default mode
		return NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil)
	}
	header := readHeaderFromFile(headerFile)
	return NewRequest(url, respType, "", "GET", "", header, nil, nil, nil)
}

func readHeaderFromFile(headerFile string) http.Header {
	//read file , parse the header and cookies
	b, err := ioutil.ReadFile(headerFile)
	if err != nil {
		//make be:  share access error
		return nil
	}
	js, _ := simplejson.NewJson(b)
	header := make(http.Header)
	jsonMap, _ := js.Map()
	for key, value := range jsonMap {
		header.Add(key, toString(value))
	}
	header.Add("Cache-Control", "max-age=0")
	header.Add("Connection", "keep-alive")
	return header
}

//point to a json file
func (this *Request) AddHeaderFile(headerFile string) *Request {
	_, err := os.Stat(headerFile)
	if err != nil {
		return this
	}
	header := readHeaderFromFile(headerFile)
	this.header = header
	return this
}

// @host  http://localhost:8765/
func (this *Request) AddProxyHost(host string) *Request {
	this.proxyHost = host
	return this
}

func (this *Request) Url() string {
	return this.url
}

func (this *Request) UrlTag() string {
	return this.urlTag
}

func (this *Request) Method() string {
	return this.method
}

func (this *Request) Postdata() string {
	return this.postData
}

func (this *Request) Header() http.Header {
	return this.header
}

func (this *Request) Cookies() []*http.Cookie {
	return this.cookies
}

func (this *Request) ProxyHost() string {
	return this.proxyHost
}

func (this *Request) ResponceType() string {
	return this.respType
}

func (this *Request) RedirectFunc() func(req *http.Request, via []*http.Request) error {

	return this.checkRedirect
}

func (this *Request) Meta() interface{} {
	return this.meta
}

func toString(v interface{}) string {
	switch f := v.(type) {
	case bool:
		if f {
			return "true"
		} else {
			return "false"
		}
	case float32:
		return strconv.FormatFloat(float64(f), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(f, 'E', -1, 64)
	case int:
		return strconv.Itoa(f)
	case int8:
		return strconv.FormatInt(int64(f), 10)
	case int16:
		return strconv.FormatInt(int64(f), 10)
	case int32:
		return strconv.FormatInt(int64(f), 10)
	case int64:
		return strconv.FormatInt(f, 10)
	case uint:
		return strconv.FormatUint(uint64(f), 10)
	case uint8:
		return strconv.FormatUint(uint64(f), 10)
	case uint16:
		return strconv.FormatUint(uint64(f), 10)
	case uint32:
		return strconv.FormatUint(uint64(f), 10)
	case uint64:
		return strconv.FormatUint(f, 10)
	case time.Time:
		return f.Format("2006-01-02 15:04:05")
	case string:
		return f
	default:
		return ""
	}
}
