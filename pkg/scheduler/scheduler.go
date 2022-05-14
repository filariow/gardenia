package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Scheduler interface {
	AddJob(JobDefinition) int64
	RemoveJob(int64)
	Start() error
	Stop() error
	TerminateRunningJobs() error
	Close()
}

type scheduler struct {
	jobIdCounter int64
	regIdCounter int64

	registrations  map[int64]jobRegistration
	mux            sync.RWMutex
	jobsToSchedule map[int64]*job
	jobsRunning    map[int64]*job
	jobsTerminated map[int64]*job

	addedJobToSchedule chan *job
	refreshSchedule    chan struct{}

	startMux  sync.Mutex
	isRunning bool
	cancel    func()
}

func New() Scheduler {
	s := scheduler{
		registrations:      make(map[int64]jobRegistration),
		jobsToSchedule:     make(map[int64]*job),
		jobsRunning:        make(map[int64]*job),
		jobsTerminated:     make(map[int64]*job),
		addedJobToSchedule: make(chan *job),
		refreshSchedule:    make(chan struct{}),
	}
	return &s
}

func (s *scheduler) Close() {
	close(s.addedJobToSchedule)
	close(s.refreshSchedule)
}

func (s *scheduler) AddJob(jd JobDefinition) int64 {
	s.mux.Lock()
	defer s.mux.Unlock()

	r := jobRegistration{
		id:  atomic.AddInt64(&s.regIdCounter, 1),
		def: jd,
	}
	s.registrations[r.id] = r
	if s.isRunning {
		s.refreshSchedule <- struct{}{}
	}
	return r.id
}

func (s *scheduler) RemoveJob(id int64) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.registrations, id)

	for k, v := range s.jobsToSchedule {
		if v.registration.id == id {
			delete(s.jobsToSchedule, k)
		}
	}
	if s.isRunning {
		s.refreshSchedule <- struct{}{}
	}
}

func (s *scheduler) TerminateRunningJobs() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, j := range s.jobsRunning {
		s.terminateJob(j)
	}
	return nil
}

func (s *scheduler) runJob(j *job) {
	s.mux.Lock()
	delete(s.jobsToSchedule, j.id)
	s.jobsRunning[j.id] = j
	s.mux.Unlock()

	if err := j.run(); err != nil {
		log.Println(err)
	}

	go func() {
		<-j.err

		s.mux.Lock()
		defer s.mux.Unlock()

		delete(s.jobsRunning, j.id)
		s.jobsTerminated[j.id] = j

		s.scheduleJob(j.registration)
	}()
}

func (s *scheduler) scheduleJob(jr jobRegistration) {
	nj := job{
		id:           atomic.AddInt64(&s.jobIdCounter, 1),
		registration: jr,
	}
	s.jobsToSchedule[nj.id] = &nj
}

func (s *scheduler) terminateJob(j *job) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.jobsRunning, j.id)
	j.cancel()
	s.jobsTerminated[j.id] = j
}

func (s *scheduler) Start() error {
	err := func() error {
		s.startMux.Lock()
		defer s.startMux.Unlock()

		if s.isRunning {
			return fmt.Errorf("Scheduler yet running")
		}
		s.isRunning = true
		return nil
	}()
	if err != nil {
		return err
	}

	// register jobs
	s.mux.Lock()
	for _, r := range s.registrations {
		s.scheduleJob(r)
	}
	s.mux.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	for {
		log.Println("waiting for new job")
		jc := s.nextJob(ctx)
		select {
		case j := <-jc:
			log.Printf("new job selected %v", j.id)
			//timer
			d := j.definition().NextSchedule().Sub(time.Now())
			timer := time.NewTimer(d)
			select {
			case <-timer.C:
				log.Println("running job")
				s.runJob(j)
				log.Println("job runned")
				log.Println("timer elapsed")
				continue
			case <-ctx.Done():
				log.Println("scheduler stopped")
				timer.Stop()
				return nil
			}
		case <-ctx.Done():
			log.Println("scheduler stopped")
			return nil
		}
	}
}

func (s *scheduler) Stop() error {
	s.startMux.Lock()
	defer s.startMux.Unlock()

	if !s.isRunning {
		return fmt.Errorf("Scheduler not running")
	}
	s.isRunning = false

	s.cancel()
	s.jobsToSchedule = map[int64]*job{}
	return nil
}

func (s *scheduler) nextJob(ctx context.Context) <-chan *job {
	c := make(chan *job)
	go func() {
		defer close(c)
		for {
			for len(s.jobsToSchedule) == 0 {
				select {
				case <-s.addedJobToSchedule:
					break
				case <-ctx.Done():
					return
				}
			}

			jb := func() *job {
				s.mux.RLock()
				defer s.mux.RUnlock()

				// find first valid job to schedul
				now := time.Now()
				var j *job
				for _, v := range s.jobsToSchedule {
					if ns := v.definition().NextSchedule(); now.Before(ns) {
						j = v
						break
					}
				}

				// if no valid job found loop again
				if j == nil {
					return nil
				}

				// look for a better candidate
				for _, v := range s.jobsToSchedule {
					ns := v.definition().NextSchedule()
					os := j.definition().NextSchedule()

					if ns.Before(os) && ns.After(now) {
						j = v
					}
				}
				return j
			}()
			if jb == nil {
				continue
			}

			// return best candidate
			c <- jb
			return
		}
	}()
	return c
}

type jobRegistration struct {
	id  int64
	def JobDefinition
}
