package main

import (
	"fmt"
	"github.com/fusuwei/gspider/gspider"
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/logger"
	"github.com/fusuwei/gspider/pkg/middleware"
	"github.com/fusuwei/gspider/pkg/rabbitmq"
	"github.com/fusuwei/gspider/pkg/request"
	"github.com/fusuwei/gspider/pkg/response"
	"net/url"
	"time"
)

type Test struct {
}

func (b *Test) Producer(f func(request *request.Request, nextNode constant.Node)) {
	for i := 0; i < 1000; i++ {
		u, err := url.Parse("http://127.0.0.1:30201/get")
		if err != nil {
			continue
		}
		r := &request.Request{
			Url:      u,
			Method:   "GET",
			Data:     map[string]interface{}{"name": "fsw"},
			Json:     nil,
			Callback: "parse",
			Meta:     nil,
			Retry:    0,
		}
		f(r, constant.ProducerNode)
	}
}

func (b *Test) RegisterCallback() map[string]gspider.Callback {
	return map[string]gspider.Callback{
		"callback": b.Callback,
		"parse":    b.Parse,
	}
}

func (b *Test) Parse(response *response.Response) (interface{}, constant.Node) {
	fmt.Println("Parse----------------")
	fmt.Println(response.Text)
	u, _ := url.Parse("http://127.0.0.1:30201/post")
	r := &request.Request{
		Url:      u,
		Method:   "POST",
		Data:     map[string]interface{}{"name": "fsw"},
		Json:     nil,
		Callback: "callback",
		Meta:     nil,
		Retry:    0,
	}
	fmt.Println("Parse----------------")
	time.Sleep(time.Second * 3)
	return r, constant.RequestNode
}

func (b *Test) Callback(response *response.Response) (interface{}, constant.Node) {
	fmt.Println("Callback----------------")
	fmt.Println(response.Text)
	fmt.Println("Callback----------------")
	return nil, constant.RequestNode
}

func main() {

	g := gspider.GSpider{
		RequestSetting: &request.Setting{
			Timeout:        30,
			Verify:         false,
			AllowRedirects: false,
		},
		GSize:         1,
		AllowCode:     []int{200, 400},
		DownloadDelay: 1,
		Model:         gspider.ConsumerModel,
		//Model:  gspider.ProducerModel,
		Logger: logger.NewLogger("debug", "test", true, "log"),
	}
	q := rabbitmq.NewQueue(
		"admin",
		"123456",
		"127.0.0.1",
		"",
		5672,
		"test")
	g.Use(constant.RequestNode, middleware.UA)
	g.Run(&Test{}, q, nil)

}
