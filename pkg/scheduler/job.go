package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type JobDefinition interface {
	NextSchedule() time.Time
	Duration() time.Duration
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

	go func() {
		j.isRun = true
		c, cancel := context.WithCancel(context.Background())
		j.cancel = cancel
		j.err = make(chan error)
		a := j.definition().Action()
		if err := a(c); err != nil {
			j.err <- err
		}
		close(j.err)
	}()

	return nil
}

func (j *job) definition() JobDefinition {
	return j.registration.def
}

type ScheduledAction func(context.Context) error
