package scheduler

import (
	"context"
	"testing"
	"time"
)

type jdempty struct {
	id int64
	t  *testing.T
}

func (jd jdempty) NextSchedule() time.Time {
	return time.Now().Add(2 * time.Second)
}

func (jd jdempty) Action() ScheduledAction {
	return func(ctx context.Context) error {
		jd.t.Logf("Empty Job %v executed", jd.id)
		return nil
	}
}

func Test_JobRun(t *testing.T) {
	t.Logf("running test job run")
	jr := jobRegistration{
		id:  1,
		def: jdempty{id: 1, t: t},
	}

	j := &job{
		id:           1,
		registration: jr,
	}
	err := j.run()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("job started with success")
	err, ok := <-j.err
	if ok {
		t.Fatal("err-chan expected to be closed")
	}

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("test completed")
}
