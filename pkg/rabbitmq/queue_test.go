package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	queue := NewQueue("admin", "eyJrZXlJZCI6ImFtcXAtdzRteDd4emtuZXEyIiwiYWxnIjoiSFMyNTYifQ.eyJzdWIiOiJhbXFwLXc0bXg3eHprbmVxMl9hZG1pbiJ9.C2aZz2xIAENZRhg3PRTPqsNIF0zB2865NWOQUhJOQvQ",
		"amqp-w4mx7xzkneq2.rabbitmq.ap-gz.public.tencenttdmq.com", "amqp-w4mx7xzkneq2|gspider", 5672, "test", 10)
	//
	//queue := NewQueue("admin", "123456", "127.0.0.1", "", 5672, "test", 10)

	queue.Start()

	//go func() {
	//	for i := 0; i < 100; i++ {
	//		msg := fmt.Sprintf("queue msg count: %d", i)
	//		fmt.Println("开始生产: ", msg)
	//		err := queue.Producer(msg)
	//		if err != nil {
	//			fmt.Println("Producer", err)
	//			return
	//		}
	//	}
	//}()

	//go func() {
	for {
		go func() {
			msg, err := queue.Consumer()
			defer func() {
				if err != nil {
					queue.Ask(msg, false)
				} else {
					queue.Ask(msg, true)
				}
			}()
			if err != nil {
				fmt.Println("Consumer", err)
				return
			}
			fmt.Println("消费:", string(msg.Msg))
			time.Sleep(time.Second * 2)
		}()

	}
	//}()

	//time.Sleep(time.Second * 100)
}
