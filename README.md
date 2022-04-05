### GSpider

GSpider是一个分布式爬虫框架，spider注重于单个网站爬取指定信息。

##### 节点

GSpider是通过节点来进行调度，共5个节点：

```go
EmptyNode    Node = iota // 空节点
ProducerNode             // 生产节点，用于生产到rabbitMq和其他队列中（选用这个节点必须实现接口queue）
RequestNode              // 请求节点
CallbackNode             // 回调节点，主要用于解析请求返回的数据
SaveNode                 // 存储节点
```

##### 队列：

只要实现以下方法可自定义队列：

```go
type Queue interface {
	Producer(msg string) error
	Consumer() (*Message, error)
	Ask(message *Message, ask bool)
	Start(size int) error
	Close()
}
```

##### 中间件：

每个节点都可以自定义中间件，中间件必须调用Next才能进入到下个中间件，Abort中止当前任务。通过GSpider.Use()来注册中间件

```go
func UA(ctx *gspider.Context) {
	ctx.Session.UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36"
	ctx.Next()
	//ctx.Abort()
}
```

示例：

> 无队列：[notQueue](https://github.com/FuSuwei/Gspider/tree/master/examples/notQueue)

> rabbitmq: [rabbitmq](https://github.com/FuSuwei/Gspider/tree/master/examples/rabbitmq)

