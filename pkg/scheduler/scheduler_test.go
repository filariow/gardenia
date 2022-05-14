package scheduler

import (
	"testing"
	"time"
)

func Test_Scheduler(t *testing.T) {
	s := (New()).(*scheduler)

	jd := jdempty{id: 1, t: t}
	i := s.AddJob(jd)

	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.jobsTerminated[i]; !ok {
		t.Fatalf("job %v not terminated", i)
	}

	if _, ok := s.jobsRunning[i]; ok {
		t.Fatalf("job %v both in running and terminated", i)
	}

	if _, ok := s.jobsToSchedule[i]; ok {
		t.Fatalf("job %v both in running and to schedule", i)
	}
}
