package workerpool

import (
	"context"
	"sync"
)

type Worker struct {
	id    int
	jobCh chan *Job
}

func NewWorker(channel chan *Job, id int) *Worker {
	return &Worker{
		id:    id,
		jobCh: channel,
	}
}

// Start starts reading a job channel while it's open
// When return is occurred wg.Done() is called
// It listens a context done channel, after signal is received the function will be return
func (wr *Worker) Start(ctx context.Context, wg *sync.WaitGroup, resCh chan<- interface{}) {
	defer wg.Done()
	for {
		select {
		case j, ok := <-wr.jobCh:
			if !ok {
				return
			}
			resCh <- j.Process(wr.id)
		case <-ctx.Done():
			return
		}
	}
}
