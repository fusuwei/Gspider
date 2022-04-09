package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

type testProducer struct {
	baseQueue
	msgs chan *Msg
}

func (t *testProducer) Producer() chan *Msg {
	return t.msgs
}

type testReceiver struct {
	baseQueue
}

func (t *testReceiver) Receiver(msg []byte) bool {
	fmt.Println(string(msg))
	return true
}

func (t *testReceiver) Qos() int {
	return 1
}

func TestRabbitMQ_Start(t *testing.T) {
	test := &testProducer{msgs: make(chan *Msg, 100)}
	defer close(test.msgs)
	//r := New("admin", "eyJrZXlJZCI6ImFtcXAtdzRteDd4emtuZXEyIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJhbXFwLXc0bXg3eHprbmVxMl9hZG1pbiJ9.C2aZz2xIAENZRhg3PRTPqsNIF0zB2865NWOQUhJOQvQ",
	//	"amqp-w4mx7xzkneq2.rabbitmq.ap-gz.public.tencenttdmq.com", "amqp-w4mx7xzkneq2|gspider", 5672)
	r := New("admin", "123456",
		"127.0.0.1", "", 5672)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			test.msgs <- &Msg{
				Msg: fmt.Sprintf("count %d", i),
			}
		}
	}()
	r.RegisterProducer(test)
	//r.RegisterReceiver(&testReceiver{})
	err := r.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.Start(10)
	time.Sleep(time.Second * 30)
}

func TestRabbitMQ_Routing(t *testing.T) {
	test := &testProducerRouting{msgs: make(chan *Msg, 100)}
	defer close(test.msgs)
	//r := New("admin", "eyJrZXlJZCI6ImFtcXAtdzRteDd4emtuZXEyIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJhbXFwLXc0bXg3eHprbmVxMl9hZG1pbiJ9.C2aZz2xIAENZRhg3PRTPqsNIF0zB2865NWOQUhJOQvQ",
	//	"amqp-w4mx7xzkneq2.rabbitmq.ap-gz.public.tencenttdmq.com", "amqp-w4mx7xzkneq2|gspider", 5672)
	r := New("admin", "123456",
		"127.0.0.1", "", 5672)

	go func() {
		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				test.msgs <- &Msg{
					RouterKey: "info",
					Msg:       fmt.Sprintf("info: %d", i),
				}
			} else {
				test.msgs <- &Msg{
					RouterKey: "error",
					Msg:       fmt.Sprintf("error: %d", i),
				}
			}
		}
	}()
	r.RegisterProducer(test)
	r.RegisterReceiver(&testReceiverRouting{})
	r.RegisterReceiver(&testReceiverRoutingError{})
	err := r.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.Start(10)
	time.Sleep(time.Second * 30)
}

type testReceiverRouting struct {
	baseQueue
}

func (t *testReceiverRouting) Qos() int {
	return 1
}

func (t *testReceiverRouting) Receiver(msg []byte) bool {
	time.Sleep(time.Second * 1)
	fmt.Println("testReceiverRouting info", string(msg))
	return true
}

func (t *testReceiverRouting) QueueName() string {
	return "routing"
}

func (t *testReceiverRouting) RouterKey() []string {
	return []string{"info"}
}

func (t *testReceiverRouting) ExchangeName() string {
	return "direct_logs"
}

func (t *testReceiverRouting) ExchangeType() string {
	return "direct"
}

func (t *testReceiverRouting) WorkType() WorkType {
	return Routing
}

type testReceiverRoutingError struct {
	baseQueue
}

func (t *testReceiverRoutingError) Qos() int {
	return 1
}

func (t *testReceiverRoutingError) Receiver(msg []byte) bool {
	time.Sleep(time.Second * 1)
	fmt.Println("testReceiverRouting error", string(msg))
	return true
}

func (t *testReceiverRoutingError) QueueName() string {
	return "routing_error"
}

func (t *testReceiverRoutingError) RouterKey() []string {
	return []string{"error"}
}

func (t *testReceiverRoutingError) ExchangeName() string {
	return "direct_logs"
}

func (t *testReceiverRoutingError) ExchangeType() string {
	return "direct"
}

func (t *testReceiverRoutingError) WorkType() WorkType {
	return Routing
}

type testProducerRouting struct {
	baseQueue
	msgs chan *Msg
}

func (t *testProducerRouting) Producer() chan *Msg {
	return t.msgs
}

func (t *testProducerRouting) QueueName() string {
	return "routing"
}

func (t *testProducerRouting) RouterKey() []string {
	return []string{"info", "waring", "error"}
}

func (t *testProducerRouting) ExchangeName() string {
	return "direct_logs"
}

func (t *testProducerRouting) ExchangeType() string {
	return "direct"
}

func (t *testProducerRouting) WorkType() WorkType {
	return Routing
}
