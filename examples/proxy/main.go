package main

import (
	"fmt"
	"github.com/fusuwei/gspider/gspider"
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/logger"
	"github.com/fusuwei/gspider/pkg/middleware"
	"github.com/fusuwei/gspider/pkg/request"
	"github.com/fusuwei/gspider/pkg/response"
	"net/url"
)

type Test struct {
}

func (b *Test) Producer(f func(request *request.Request, nextNode constant.Node)) {
	for i := 0; i < 2; i++ {
		u, err := url.Parse("http://47.118.61.27:30201/get")
		if err != nil {
			continue
		}
		r := &request.Request{
			Url:      u,
			Method:   "GET",
			Json:     map[string]interface{}{"name": "fsw"},
			Data:     nil,
			Callback: "parse",
			Meta:     nil,
			Retry:    0,
		}
		f(r, constant.RequestNode)
	}
}

func (b *Test) RegisterCallback() map[string]gspider.Callback {
	return map[string]gspider.Callback{
		"parse": b.Parse,
	}
}

func (b *Test) Parse(response *response.Response) (interface{}, constant.Node) {
	fmt.Println("Parse----------------")
	fmt.Println(response.Text)
	fmt.Println("Parse----------------")
	return nil, constant.EmptyNode
}

func main() {
	g := gspider.Default()
	l := logger.NewLogger("debug", "test", true, "log")
	g.Logger = l
	//m := mysql.NewMysql("root", "123456", "127.0.0.1", "Test", 3306, l)
	//m.OpenConnection()
	//T := TestSaver{mysql: m}
	g.Use(constant.RequestNode, middleware.UA)
	g.Use(constant.RequestNode, middleware.SetProxy)
	g.Run(&Test{}, nil, nil)
}
