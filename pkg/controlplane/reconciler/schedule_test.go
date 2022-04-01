package reconciler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/controlplane"
	"github.com/barrydevp/transcoorditor/pkg/controlplane/reconciler"
)

type entry struct {
	e  time.Time
	id string
}

func (e *entry) ExpiredAt() *time.Time {
	return &e.e
}

func TestScheduleReconciler(t *testing.T) {
	returns := []reconciler.ScheduleEntry{
		&entry{
			e:  time.Now().Add(0 * time.Second),
			id: "abd",
		},
		&entry{
			e:  time.Now().Add(0 * time.Second),
			id: "abd-1",
		},
		&entry{
			e:  time.Now().Add(1 * time.Second),
			id: "abdc",
		},
	}
	recl := reconciler.NewScheduleReconciler(
		func() []reconciler.ScheduleEntry {
			return nil
		},
		func(entries []reconciler.ScheduleEntry) []reconciler.ScheduleEntry {
			for _, en := range entries {
				fmt.Println(en.(*entry).id)
			}
			rets := returns
			returns = []reconciler.ScheduleEntry{}

			return rets
		})
	ctrlplane := controlplane.New(nil)

	ctrlplane.RegisterRecl(recl)

	ctrlplane.Run()

	now := time.Now()

	testCases := []struct {
		d  time.Duration
		id string
	}{
		{
			d:  0,
			id: "s1",
		},
		{
			d:  0,
			id: "s2",
		},
		{
			d:  0,
			id: "s3",
		},
		{
			d:  3,
			id: "s4",
		},
		{
			d:  4,
			id: "s5",
		},
	}

	for _, test := range testCases {
		recl.Schedule(&entry{
			e:  now.Add(test.d * time.Second),
			id: test.id,
		})
	}

	time.Sleep(1 * time.Second)
	now = time.Now()

	for _, test := range testCases {
		recl.Schedule(&entry{
			e:  now.Add(test.d * time.Second),
			id: test.id,
		})
	}

	time.Sleep(15 * time.Second)

	ctrlplane.Stop()
}
