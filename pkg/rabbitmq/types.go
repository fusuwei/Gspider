package rabbitmq

import "fmt"

type base interface {
	QueueName() string   // 获取接收者需要监听的队列
	RouterKey() []string // 这个队列绑定的路由
	OnError(error)       // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	ExchangeName() string
	ExchangeType() string
	WorkType() WorkType
}

type Msg struct {
	RouterKey string
	Msg       string
}

// Producer 定义生产者接口
type Producer interface {
	Producer() chan *Msg
	base
}

// Receiver 定义接收者接口
type Receiver interface {
	Receiver([]byte) bool
	Qos() int
	base
}

type baseQueue struct{}

func (b *baseQueue) QueueName() string {
	return "gspider"
}

func (b *baseQueue) RouterKey() []string {
	return []string{""}
}

func (b *baseQueue) OnError(err error) {
	fmt.Println(err)
}

func (b *baseQueue) ExchangeName() string {
	return ""
}

func (b *baseQueue) ExchangeType() string {
	return ""
}

func (b *baseQueue) WorkType() WorkType {
	return WorkQueue
}
