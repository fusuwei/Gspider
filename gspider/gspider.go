package gspider

import (
	"errors"
	"fmt"
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/logger"
	"github.com/fusuwei/gspider/pkg/queue"
	"github.com/fusuwei/gspider/pkg/request"
	"github.com/fusuwei/gspider/pkg/response"
	"github.com/fusuwei/gspider/pkg/storage"
	"time"
)

type Model string

const (
	ProducerModel = "Producer"
	ConsumerModel = "Consumer"
)

type GSpider struct {
	AllowCode      []int
	DownloadDelay  int
	Model          Model
	GSize          int
	Callbacks      map[string]Callback
	RequestSetting *request.Setting
	KeepSession    bool

	spider Spider
	queue  queue.Queue
	saver  storage.Saver

	Logger *logger.Logger

	session  *request.Session
	workers  *Workers
	lastTime time.Time

	handlerChains map[constant.Node]HandlerChain
}

func Default() *GSpider {
	return &GSpider{
		RequestSetting: &request.Setting{
			Timeout:        30,
			Verify:         false,
			AllowRedirects: false,
		},
		AllowCode:     []int{200},
		DownloadDelay: 0,
		GSize:         1,
		Model:         ProducerModel,
		Logger:        logger.NewLogger("debug", "GSpider", false, "log"),
	}
}

func (g *GSpider) Run(spider Spider, queue queue.Queue, save storage.Saver) {
	g.init(spider, queue, save)
	g.start()
}

func (g *GSpider) init(spider Spider, queue queue.Queue, save storage.Saver) {
	if g.GSize == 0 {
		g.GSize = 1
	}
	g.workers = NewWorkers(g.GSize)

	g.spider = spider
	g.queue = queue
	g.saver = save
	g.Callbacks = g.spider.RegisterCallback()
	if g.Logger == nil {
		g.Logger = logger.NewLogger("debug", "GSpider", true, "log")
	}
	if g.RequestSetting != nil {
		var timeout time.Duration
		if g.RequestSetting.Timeout != 0 {
			timeout = time.Second * time.Duration(g.RequestSetting.Timeout)
		} else {
			timeout = time.Second * 30
		}
		g.session = request.New(timeout, g.RequestSetting.Verify, g.RequestSetting.AllowRedirects)
	} else {
		g.session = request.New(time.Second*30, false, false)
	}
}

func (g *GSpider) start() {
	if g.queue != nil {
		err := g.queue.Start(g.GSize)
		if err != nil {
			g.Logger.Errorf("初始化队列失败, error: %s", err.Error())
		}
		defer g.queue.Close()
	}
	switch g.Model {
	case constant.ProducerModel:
		go g.spider.Producer(g.submit)
	case constant.ConsumerModel:
		go g.consumer()
	default:
		go g.spider.Producer(g.submit)
		go g.consumer()
	}
	g.dispatch()
}

func (g *GSpider) submit(r *request.Request, nextNode constant.Node) {
	ctx := newContext(g.session, g.Logger)
	ctx.Request = r
	ctx.NextNode = nextNode
	ctx.HandlerChains = g.handlerChains
	g.workers.PushTask(ctx)
}

func (g *GSpider) consumer() {
	if g.queue == nil {
		g.Logger.Errorf("未初始化队列")
	}
	for {
		func() {
			msg, err := g.queue.Consumer()
			defer func() {
				if err != nil {
					g.Logger.Errorf("消费失败, error:%s", err.Error())
					g.queue.Ask(msg, false)
					return
				}
			}()
			if msg == nil {
				return
			}
			if err != nil {
				return
			}
			g.Logger.Debugf("开始消费: %s", string(msg.Msg))
			res, err := request.ToRequest(string(msg.Msg))
			if err != nil {
				return
			}
			ctx := newContext(g.session, g.Logger)
			ctx.Request = res
			ctx.NextNode = constant.RequestNode
			ctx.HandlerChains = g.handlerChains
			ctx.ask = msg.Ask
			g.workers.PushTask(ctx)
		}()
	}
}

func (g *GSpider) producer(ctx *Context) {
	if g.queue == nil {
		g.Logger.Errorf("未初始化队列，无法生产")
		g.workers.Finish(ctx, false)
		return
	}
	if ctx.Request == nil {
		g.Logger.Errorf("请求体为空")
		g.workers.Finish(ctx, false)
		return
	}
	publishData := ctx.Request.ToPublish()
	g.Logger.Debugf("开始生产：%s", publishData)
	if err := g.queue.Producer(publishData); err != nil {
		g.Logger.Errorf("生产错误，error: %s", err.Error())
		g.workers.Finish(ctx, false)
		return
	}
	g.workers.Finish(ctx, true)
}

func (g *GSpider) down(ctx *Context) {
	status := false
	ok := false
	var err error
	var res *response.Response
	defer func() {
		if !status {
			g.workers.Finish(ctx, false)
		}
	}()
	if ctx.Request == nil {
		g.Logger.Errorf("请求体为空")
		return
	}
	for i := 0; i < (ctx.Request.Retry + 1); i++ {
		if g.DownloadDelay != 0 {
			g.Logger.Debugf(fmt.Sprintf("延迟%d秒请求", g.DownloadDelay))
			time.Sleep(time.Second * time.Duration(g.DownloadDelay))
		}
		g.Logger.Debugf("开始请第%d次：%s", i+1, ctx.Request.Url.String())
		res, err = ctx.Session.Request(ctx.Request)
		if err != nil {
			g.Logger.Errorf("请求失败，error: %s", err.Error())
			continue
		}
		flag := false
		for _, statusCode := range g.AllowCode {
			if res.Resp.StatusCode == statusCode {
				flag = true
				ok = true
			}
		}
		if ok {
			break
		}
		if !flag {
			g.Logger.Warnf("状态码错误, status code: %d", res.Resp.StatusCode)
			continue
		}
	}
	if ok {
		status = true
	} else {
		status = false
		return
	}
	if g.KeepSession {
		g.session = ctx.Session
	}
	ctx.Response = res
	ctx.NextNode = constant.CallbackNode
	status = true
	ctx.indexes[ctx.NextNode] = -1
	g.workers.PushWaitTask(ctx)
}

func (g *GSpider) call(ctx *Context) {
	res := ctx.Response
	if res == nil {
		g.Logger.Errorf("响应体为空")
		g.workers.Finish(ctx, false)
		return
	}
	if call, ok := g.Callbacks[res.Callback]; !ok {
		g.Logger.Errorf("未注册回调函数，callback: %s", res.Callback)
		g.workers.Finish(ctx, false)
		return
	} else {
		result, node := call(res)
		if result == nil {
			g.workers.Finish(ctx, true)
			return
		}
		err := g.checkType(result, ctx)
		if err != nil {
			g.Logger.Errorf(err.Error())
			g.workers.Finish(ctx, false)
			return
		}
		ctx.NextNode = node
		ctx.indexes[ctx.NextNode] = -1
		g.workers.PushWaitTask(ctx)
	}
}

func (g *GSpider) save(ctx *Context) {
	if ctx.Item == nil {
		g.Logger.Errorf("存储体为空")
		return
	}
	err := g.saver.Save(ctx.Item)
	if err != nil {
		g.workers.Finish(ctx, false)
		return
	}
	if ctx.Item.NextPage != nil {
		if ctx.Item.NextPageNode == constant.EmptyNode {
			g.Logger.Warnf("next page node is empty node")
		}
		go g.submit(ctx.Item.NextPage, ctx.Item.NextPageNode)
	}
	g.workers.Finish(ctx, true)
}

func (g *GSpider) dispatch() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if g.lastTime.Unix() > 0 && time.Now().Unix()-g.lastTime.Unix() > 5*60 {
				g.Logger.Debugf("消费超时")
				return
			}
		case work := <-g.workers.GetWaitWorkers():
			g.lastTime = time.Now()
			work.handle()
			if work.isAbort {
				g.workers.Finish(work, false)
				continue
			}
			switch work.NextNode {
			case constant.ProducerNode:
				go g.producer(work)
			case constant.RequestNode:
				go g.down(work)
			case constant.CallbackNode:
				go g.call(work)
			case constant.SaveNode:
				go g.save(work)
			case constant.EmptyNode:
				g.workers.Finish(work, true)

			}
		}
	}
}

func (g *GSpider) checkType(data interface{}, ctx *Context) (err error) {
	if data == nil {
		return nil
	}
	switch data.(type) {
	case *request.Request:
		ctx.Request = data.(*request.Request)
	case *response.Response:
		ctx.Response = data.(*response.Response)
	case request.Request:
		res := data.(request.Request)
		ctx.Request = &res
	case response.Response:
		res := data.(response.Response)
		ctx.Response = &res
	case storage.Item:
		res := data.(storage.Item)
		ctx.Item = &res
	case *storage.Item:
		res := data.(*storage.Item)
		ctx.Item = res
	default:
		err = errors.New("返回类型错误！")
	}
	return
}

func (g *GSpider) Use(node constant.Node, h Handler) {
	if g.handlerChains == nil {
		g.handlerChains = make(map[constant.Node]HandlerChain)
	}
	if handlers, ok := g.handlerChains[node]; ok {
		g.handlerChains[node] = append(handlers, h)
	} else {
		g.handlerChains[node] = HandlerChain{h}
	}
}
