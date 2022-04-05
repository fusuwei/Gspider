package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type WorkType int

var (
	WorkQueue        = WorkType(1)
	PublishSubscribe = WorkType(2)
	Routing          = WorkType(3)
	Topics           = WorkType(4)
)

type RabbitMQ struct {
	User     string
	password string
	Host     string
	Port     int
	vhost    string

	Type WorkType

	qos               int
	producerList      []Producer
	receiverList      []Receiver
	connection        *amqp.Connection
	channel           *amqp.Channel
	lock              sync.RWMutex
	isRunningReceiver bool // 接收者是否启动
	cancel            context.CancelFunc
}

func New(user, pwd, host, vhost string, port int) *RabbitMQ {
	return &RabbitMQ{
		User:         user,
		password:     pwd,
		Host:         host,
		Port:         port,
		vhost:        vhost,
		producerList: make([]Producer, 0),
		receiverList: make([]Receiver, 0),
	}
}

// Connect 连接rabbitmq
func (r *RabbitMQ) Connect() (err error) {
	//logger.Debug(fmt.Sprintf("开始连接rabbitmq[%s]队列！", r.Host))
	u := "amqp://" + r.User + ":" + r.password + "@" + r.Host + ":5672/"
	if r.vhost != "" {
		u = u + r.vhost
	}

	r.connection, err = amqp.Dial(u)
	if err != nil {
		return
	}
	r.channel, err = r.connection.Channel()
	if err != nil {
		return
	}
	return nil
}

// Close 关闭RabbitMQ连接
func (r *RabbitMQ) Close() (err error) {
	// 先关闭管道,再关闭链接
	r.cancel()
	err = r.channel.Close()
	if err != nil {
		return
	}
	err = r.connection.Close()
	if err != nil {
		return
	}
	return nil
}

// RegisterReceiver 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.lock.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.lock.Unlock()
}

// RegisterProducer 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.producerList = append(r.producerList, producer)
}

// Start 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) Start(size int) {
	// 开启监听生产者发送任务
	r.isRunningReceiver = false
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.qos = size
	for _, producer := range r.producerList {
		go r.listenProducer(producer, ctx)
	}
	// 开启监听接收者接收任务
	for _, receiver := range r.receiverList {
		go r.listenReceiver(receiver, ctx)
	}
	time.Sleep(1 * time.Second)
}

// 发送任务
func (r *RabbitMQ) listenProducer(p Producer, ctx context.Context) (err error) {
	queueName := p.QueueName()
	// 用于检查队列是否存在,已经存在不需要重复声明

	if p.WorkType() == WorkQueue {
		_, err = r.channel.QueueDeclare(queueName, true, false, false, true, nil)
		if err != nil {
			p.OnError(err)
			return err
		}
	} else if p.WorkType() == Routing || p.WorkType() == PublishSubscribe {
		r.wait()
	}
	// 发送任务消息
	channel := p.Producer()
	for msg := range channel {
		select {
		case <-ctx.Done():
			return nil
		default:
			var key string
			if msg.RouterKey != "" {
				key = msg.RouterKey
			} else {
				key = p.QueueName()
			}
			err = r.channel.Publish(p.ExchangeName(), key, false, false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(msg.Msg),
				})
			p.OnError(err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver, ctx context.Context) error {
	defer r.Close()
	_, err := r.channel.QueueDeclare(receiver.QueueName(), true, false, false, true, nil)
	if err != nil {
		receiver.OnError(err)
		return err
	}

	if receiver.WorkType() == Routing || receiver.WorkType() == PublishSubscribe {
		err = r.channel.ExchangeDeclare(receiver.ExchangeName(), receiver.ExchangeType(), true, false, false, true, nil)
		if err != nil {
			err = errors.New(fmt.Sprintf("MQ注册交换机失败:%s", err))
			receiver.OnError(err)
			return err
		}

		for _, routerKey := range receiver.RouterKey() {
			err := r.channel.QueueBind(receiver.QueueName(), routerKey, receiver.ExchangeName(), true, nil)
			if err != nil {
				err = errors.New(fmt.Sprintf("MQ绑定队列失败:%s ", err))
				receiver.OnError(err)
				return err
			}
		}
	}

	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(r.qos, 0, true)
	msgList, err := r.channel.Consume(receiver.QueueName(), "", false, false, false, true, nil)
	if err != nil {
		receiver.OnError(err)
		return errors.New(fmt.Sprintf("获取消费通道异常:%s", err))
	}
	r.isRunningReceiver = true
	for msg := range msgList {
		// 处理数据
		select {
		case <-ctx.Done():
			return nil
		default:
			go func(msg amqp.Delivery) {
				if ok := receiver.Receiver(msg.Body); ok {
					err = msg.Ack(false)
				} else {
					err = msg.Ack(true)
				}
			}(msg)
		}
	}
	return nil
}

func (r *RabbitMQ) wait() bool {
	for r.isRunningReceiver {
		time.Sleep(time.Millisecond * 200)
		return true
	}
	return false
}
