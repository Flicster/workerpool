package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	Jobs []*Job

	size      int
	collector chan *Job
}

func NewPool(jobs []*Job, size int) *Pool {
	return &Pool{
		Jobs:      jobs,
		size:      size,
		collector: make(chan *Job, len(jobs)),
	}
}

// Run runs amount of workers according to the size parameter
// Those workers process passed jobs
// It listens a context done channel, after signal is received the pool awaiting jobs in processing and return results
func (p *Pool) Run(ctx context.Context) (res []interface{}, errs []error) {
	defer func() {
		for _, j := range p.Jobs {
			if j.Err != nil {
				errs = append(errs, j.Err)
			}
		}
	}()

	resCh := make(chan interface{}, len(p.Jobs))
	wg := &sync.WaitGroup{}
	wg.Add(p.size)
	for i := 1; i <= p.size; i++ {
		worker := NewWorker(p.collector, i)
		go worker.Start(ctx, wg, resCh)
	}

	for _, j := range p.Jobs {
		p.collector <- j
	}
	close(p.collector)

	wg.Wait()
	close(resCh)

	for r := range resCh {
		if r != nil {
			res = append(res, r)
		}
	}

	return
}
