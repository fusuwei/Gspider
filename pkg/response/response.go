package response

import (
	"compress/gzip"
	"github.com/fusuwei/gspider/pkg/charset"
	"github.com/fusuwei/gspider/pkg/xpath"
	"io/ioutil"
	"net/http"
	"strings"
)

type Response struct {
	Resp     *http.Response
	Content  []byte
	Text     string
	Callback string
	Meta     map[string]interface{}
	Dom      *xpath.Select
}

func New(res *http.Response, callback string, meta map[string]interface{}) *Response {
	defer res.Body.Close()
	response := &Response{
		Resp:     res,
		Callback: callback,
		Meta:     meta,
	}
	body := res.Body
	if res.Header.Get("Content-Encoding") == "gzip" && res.Header.Get("Accept-Encoding") != "" {
		r, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil
		}
		body = r
	}
	content, err := ioutil.ReadAll(body)
	if err != nil {
		response.Content = nil
		return response
	}
	response.Content = content
	response.Text = response.initBody()
	response.Dom = response.initDom()
	return response
}

func (r *Response) initBody() string {
	contentType := r.Resp.Header.Get("Content-Type")
	var c charset.Charset
	if strings.Contains(contentType, "GBK") || strings.Contains(contentType, "gbk") {
		c = charset.GBK
	} else if strings.Contains(contentType, "GB2312") || strings.Contains(contentType, "gb2312") {
		c = charset.GB2312

	} else if strings.Contains(contentType, "GB18030") || strings.Contains(contentType, "gb18030") {
		c = charset.GB18030
	} else {
		c = charset.UTF_8
	}
	text, err := charset.ToUTF8(c, string(r.Content))
	if err != nil {
		return ""
	}
	return text
}

func (r *Response) initDom() *xpath.Select {
	if r.Text == "" {
		return nil
	}
	return xpath.Parse(r.Text)
}
