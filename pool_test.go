package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestPool(t *testing.T) {
	s := new(PoolTestSuit)
	suite.Run(t, s)
}

type PoolTestSuit struct {
	suite.Suite
	jobs []*Job
}

func (p *PoolTestSuit) SetupSuite() {
	jobLen := 10
	p.jobs = make([]*Job, 0, jobLen)
	for x := 0; x < jobLen; x++ {
		x := x
		p.jobs = append(p.jobs, NewJob(func(data interface{}) (res interface{}, err error) {
			time.Sleep(time.Second * 1)
			if x == 0 {
				return nil, errors.New("test error")
			}
			return data, nil
		}, 1))
	}
	return
}

func (p *PoolTestSuit) TestRunWithSingleWorker() {
	now := time.Now()
	pool := NewPool(p.jobs, 1)
	res, err := pool.Run(context.Background())

	p.Equal(9, len(res))
	p.Equal(1, len(err))
	p.WithinDuration(now.Add(time.Second*10), time.Now(), time.Millisecond*150)
}

func (p *PoolTestSuit) TestRunWithDoubleWorker() {
	now := time.Now()
	pool := NewPool(p.jobs, 2)
	res, err := pool.Run(context.Background())

	p.Equal(9, len(res))
	p.Equal(1, len(err))
	p.WithinDuration(now.Add(time.Second*5), time.Now(), time.Millisecond*150)
}

func (p *PoolTestSuit) TestRunWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	now := time.Now()
	pool := NewPool(p.jobs, 1)
	pool.Run(ctx)

	p.WithinDuration(now.Add(time.Second*1), time.Now(), time.Second*4)
}
