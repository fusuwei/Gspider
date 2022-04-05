package storage

type Item struct {
	Id   int
	Body interface{}
}

type Saver interface {
	Save(*Item) error
}
