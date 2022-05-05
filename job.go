package workerpool

import (
	"errors"
	"fmt"
)

type JobHandler func(data interface{}) (res interface{}, err error)

type Job struct {
	Err     error
	Data    interface{}
	handler JobHandler
}

func NewJob(handler JobHandler, data interface{}) *Job {
	return &Job{
		handler: handler,
		Data:    data,
	}
}

// Process handles a job with pass a job data into it
// If error occurred it saves this error into Err field
func (j *Job) Process(workerId int) interface{} {
	res, err := j.handler(j.Data)
	if err != nil {
		j.Err = errors.New(fmt.Sprintf("worker: %d failed a job with error: %s", workerId, err))
	}
	return res
}
