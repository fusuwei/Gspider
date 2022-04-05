package queue

type Message struct {
	Msg []byte
	Ask chan bool
}

type Queue interface {
	Producer(msg string) error
	Consumer() (*Message, error)
	Ask(message *Message, ask bool)
	Start(size int) error
	Close()
}
