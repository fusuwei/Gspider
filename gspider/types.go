package gspider

import (
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/request"
	"github.com/fusuwei/gspider/pkg/response"
)

type Callback func(response *response.Response) (interface{}, constant.Node)

type Spider interface {
	Producer(func(request *request.Request, nextNode constant.Node))
	RegisterCallback() map[string]Callback
}

type Handler func(ctx *Context)

type HandlerChain []Handler
