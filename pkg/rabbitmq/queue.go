package rabbitmq

import "github.com/fusuwei/gspider/pkg/queue"

type RabbitProducer struct {
	baseQueue
	name   string
	errors chan error
	msgs   chan *Msg
}

func (r *RabbitProducer) QueueName() string {
	return r.name
}

func (r *RabbitProducer) Producer() chan *Msg {
	return r.msgs
}

func (r *RabbitProducer) OnError(err error) {
	r.errors <- err
}

type RabbitReceiver struct {
	msg    chan *queue.Message
	errors chan error
	qos    int
	name   string
	baseQueue
}

func (r *RabbitReceiver) Receiver(msg []byte) bool {
	m := &queue.Message{
		Msg: msg,
		Ask: make(chan bool),
	}
	r.msg <- m
	return <-m.Ask
}

func (r *RabbitReceiver) Qos() int {
	return r.qos
}

func (r *RabbitReceiver) OnError(err error) {
	r.errors <- err
}

func (r *RabbitReceiver) QueueName() string {
	return r.name
}

type Queue struct {
	*RabbitProducer
	*RabbitReceiver
	*RabbitMQ
}

func NewQueue(user, pwd, host, vhost string, port int, name string) *Queue {
	rabbitProducer := &RabbitProducer{
		name:   name,
		msgs:   make(chan *Msg),
		errors: make(chan error),
	}
	rabbitReceiver := &RabbitReceiver{
		msg:  make(chan *queue.Message,),
		name: name,
	}
	rabbitmq := New(user, pwd, host, vhost, port)
	rabbitmq.RegisterProducer(rabbitProducer)
	rabbitmq.RegisterReceiver(rabbitReceiver)
	return &Queue{
		rabbitProducer,
		rabbitReceiver,
		rabbitmq,
	}
}

func (r *Queue) Producer(msg string) error {
	select {
	case err := <-r.RabbitProducer.errors:
		return err
	default:
		r.RabbitProducer.msgs <- &Msg{
			RouterKey: "",
			Msg:       msg,
		}
		return <-r.RabbitProducer.errors
	}
}

func (r *Queue) Consumer() (*queue.Message, error) {
	select {
	case err := <-r.RabbitReceiver.errors:
		return nil, err
	default:
		msg := <-r.RabbitReceiver.msg
		return msg, nil
	}
}

func (r *Queue) Ask(message *queue.Message, ask bool) {
	message.Ask <- ask
}

func (r *Queue) Start(size int) error {
	//if reStart{
	//	rabbitmq := New(r.RabbitMQ.User, r.RabbitMQ.password, r.RabbitMQ.Host, r.RabbitMQ.vhost, r.RabbitMQ.Port)
	//	rabbitmq.RegisterProducer(r.RabbitProducer)
	//	rabbitmq.RegisterReceiver(r.RabbitReceiver)
	//}
	err := r.RabbitMQ.Connect()
	if err != nil {
		return err
	}
	r.RabbitMQ.Start(size)
	return nil
}

func (r *Queue) Close() {
	r.RabbitMQ.Close()
	close(r.RabbitReceiver.msg)
	close(r.RabbitProducer.msgs)
}
