package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type JobDefinition interface {
	NextSchedule() time.Time
	Action() ScheduledAction
}

type job struct {
	id           int64
	registration jobRegistration

	cancel func()
	err    chan error

	mux   sync.Mutex
	isRun bool
}

func (j *job) run() error {
	j.mux.Lock()
	defer j.mux.Unlock()

	if j.isRun {
		return fmt.Errorf("Job yet run")
	}

	j.err = make(chan error)
	c, cancel := context.WithCancel(context.Background())
	j.cancel = cancel

	go func() {
		j.isRun = true
		defer close(j.err)
		a := j.definition().Action()
		if err := a(c); err != nil {
			j.err <- err
		}
	}()

	return nil
}

func (j *job) definition() JobDefinition {
	return j.registration.def
}

type ScheduledAction func(context.Context) error
