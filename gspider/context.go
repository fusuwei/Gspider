package gspider

import (
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/logger"
	"github.com/fusuwei/gspider/pkg/request"
	"github.com/fusuwei/gspider/pkg/response"
	"github.com/fusuwei/gspider/pkg/storage"
)

type Context struct {
	Proxy         map[string]string
	Session       *request.Session
	Request       *request.Request
	Response      *response.Response
	Item          *storage.Item
	indexes       map[constant.Node]int
	logger        *logger.Logger
	HandlerChains map[constant.Node]HandlerChain
	NextNode      constant.Node

	ask     chan bool
	isAbort bool
}

func (c *Context) Next() {
	if _, ok := c.indexes[c.NextNode]; !ok {
		c.indexes[c.NextNode] = -1
	}
	if _, ok := c.HandlerChains[c.NextNode]; !ok {
		return
	}
	c.indexes[c.NextNode]++
	for c.indexes[c.NextNode] < len(c.HandlerChains[c.NextNode]) {
		c.HandlerChains[c.NextNode][c.indexes[c.NextNode]](c)
		c.indexes[c.NextNode]++
	}
}

func (c *Context) Abort() {
	c.logger.Errorf("中间件中止任务")
	index := c.indexes[c.NextNode]
	index = len(c.HandlerChains[c.NextNode]) + 1
	c.indexes[c.NextNode] = index
	c.isAbort = true
}

func newContext(session *request.Session, log *logger.Logger) *Context {
	return &Context{
		Session:  session,
		NextNode: constant.EmptyNode,
		logger:   log,
		indexes:  make(map[constant.Node]int),
	}
}

func (c *Context) done(ask bool) {
	if c.ask == nil {
		return
	}
	c.ask <- ask
}

func (c *Context) handle() {
	if c.NextNode == constant.EmptyNode {
		return
	}
	c.Next()
}
