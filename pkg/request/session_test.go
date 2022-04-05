package request

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestSession_Request(t *testing.T) {
	session := New("User-Agent", time.Second*30, false, false, map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36",
	})
	u, _ := url.Parse("https://ip.jiangxianli.com/?page=1")
	r := &Request{
		Url:      u,
		Method:   "GET",
		Data:     nil,
		Json:     nil,
		Callback: "",
		Meta:     nil,
		Retry:    0,
	}
	res, err := session.Request(r)
	if err != nil {
		t.Log(err.Error())
		return
	}
	trs := res.Dom.Xpath("//table[@class='layui-table']/tbody/tr")
	for _, tr := range trs.NextDocs {
		ip := tr.Xpath("./td[1]").ExtractFirst()
		fmt.Println(ip)
	}
}
