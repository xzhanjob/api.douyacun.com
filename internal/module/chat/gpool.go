package chat

import (
	"fmt"
	"time"
)

type pool struct {
	sem  chan struct{}
	work chan func()
}

var ErrTimeout = fmt.Errorf("schedule error: timed out")

func NewPool(size, queue int) *pool {
	return &pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}
}
func (p *pool) Schedule(task func()) error {
	return p.schedule(task, nil)
}

func (p *pool) ScheduleTimeout(task func(), timeout time.Duration) error {
	return p.schedule(task, time.After(timeout))
}

func (p *pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}

func (p *pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()
	task()
	for task := range p.work {
		task()
	}
}
