package gspider

type Workers struct {
	freeWorkers chan struct{}
	waitWorkers chan *Context
}

func NewWorkers(size int) *Workers {
	workers := &Workers{
		waitWorkers: make(chan *Context),
		freeWorkers: make(chan struct{}, size),
	}
	return workers
}

func (w *Workers) PushTask(context *Context) {
	w.freeWorkers <- struct{}{}
	w.waitWorkers <- context
}

func (w *Workers) PushWaitTask(context *Context) {
	w.waitWorkers <- context
}

func (w *Workers) Close() {
	close(w.freeWorkers)
	close(w.waitWorkers)
}

func (w *Workers) Finish(worker *Context, result bool) {
	worker.done(result)
	<-w.freeWorkers
}

func (w *Workers) GetWaitWorkers() <-chan *Context {
	return w.waitWorkers
}
