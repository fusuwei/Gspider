package storage

import (
	"github.com/fusuwei/gspider/pkg/constant"
	"github.com/fusuwei/gspider/pkg/request"
)

type Item struct {
	Id           int
	Body         interface{}
	NextPage     *request.Request
	NextPageNode constant.Node
}

type Saver interface {
	Save(*Item) error
}
