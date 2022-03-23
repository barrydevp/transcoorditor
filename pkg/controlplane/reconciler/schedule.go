package reconciler

import (
	"container/heap"
	"sync"
	"time"
)

type ScheduleEntry interface {
	ExpiredAt() *time.Time
}

type entryPQ []ScheduleEntry

func (pq entryPQ) Len() int { return len(pq) }

func (pq entryPQ) Less(i, j int) bool {
	// We want Pop to give us the nearlest, not farest, priority so we use greater than here.
	return pq[i].ExpiredAt().Before(*pq[j].ExpiredAt())
}

func (pq entryPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *entryPQ) Push(x interface{}) {
	item := x.(ScheduleEntry)
	*pq = append(*pq, item)
}

func (pq *entryPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (pq *entryPQ) Peek() ScheduleEntry {
	if len(*pq) > 0 {
		return (*pq)[0]
	}
	return nil
}

type HandleExpiredFunc func(entries []ScheduleEntry) []ScheduleEntry
type InitScheQueueFunc func() []ScheduleEntry

type ScheduleReconciler struct {
	handleExpired HandleExpiredFunc
	initScheQueue InitScheQueueFunc

	mutex       sync.Mutex
	scheQueue   *entryPQ
	waiting     bool
	interruptCh *chan struct{}
}

func NewScheduleReconciler(initScheQueue InitScheQueueFunc, handleExpired HandleExpiredFunc) *ScheduleReconciler {
	return &ScheduleReconciler{
		handleExpired: handleExpired,
		initScheQueue: initScheQueue,
		scheQueue:     &entryPQ{},
	}
}

// returned value for determining whether or not we should skip enter into waiting phase
// and continue the loop due to new entries has been added
func (r *ScheduleReconciler) handleExpiredEntries(entries []ScheduleEntry) bool {
	if r.handleExpired != nil {
		newEntries := r.handleExpired(entries)
		if len(newEntries) > 0 {
			r.ScheduleBatch(newEntries)
			logger.Info("new entries added: ", len(newEntries))

			// skip waiting, continue the loop
			return true
		}
	} else {
		logger.Warn("unhandled expired retries")
	}

	return false
}

func (r *ScheduleReconciler) getExpiredEntriesAndNext(now time.Time) ([]ScheduleEntry, ScheduleEntry) {
	var expireds []ScheduleEntry

	for r.scheQueue.Peek() != nil {
		en := r.scheQueue.Peek()
		if now.Before(*en.ExpiredAt()) {
			return expireds, en
		}

		expireds = append(expireds, heap.Pop(r.scheQueue).(ScheduleEntry))
	}

	return expireds, nil
}

// this must be called once at a time (mean you must Lock the mutex)
func (r *ScheduleReconciler) interupt() {
	// if reconciler is waiting, signal it and set watiting to false
	if r.interruptCh != nil && r.waiting {
		*r.interruptCh <- struct{}{}
		r.waiting = false
	}
}

func (r *ScheduleReconciler) Bootstrap() {
	entries := r.initScheQueue()
	r.ScheduleBatch(entries)
}

// this should running inside an goroutine and call once at a time
func (r *ScheduleReconciler) Reconcile(stopCh <-chan struct{}) {
	// r.mutex.Lock()
	if r.interruptCh != nil {
		// r.mutex.Unlock()
		return
	}
	// r.mutex.Unlock()

	logger.Info("reconciling...")

	interCh := make(chan struct{})
	r.interruptCh = &interCh

	for {
		now := time.Now()

		// handle/check phase
		r.mutex.Lock()
		expireds, next := r.getExpiredEntriesAndNext(now)
		logger.Info(len(expireds), " expired entries")

		if len(expireds) > 0 {
			r.mutex.Unlock()
			// handle expired entries
			_ = r.handleExpiredEntries(expireds) // skipWait

			// when handleExpiredEntries running maybe new job has been added
			// So we must recheck the queue
			continue
		}

		r.waiting = true
		r.mutex.Unlock()

		// waiting phase
		var timer *time.Timer

		if next != nil {
			timer = time.NewTimer(next.ExpiredAt().Sub(now))
			logger.Info("now: ", now, "next: ", next.ExpiredAt)
		} else {
			// no entry need to schedule yet, sleep until new entry was added
			timer = time.NewTimer(100000 * time.Hour)
			logger.Info("sleep until new entry was added")
		}

		select {
		case <-timer.C:
			// next expired
		case <-*r.interruptCh:
			// usually a new entry has added
			logger.Info("received interrupt")
		case <-stopCh:
			// received stop sigal
			// @TODO cleanup code
			break
		}
	}
}

func (r *ScheduleReconciler) Schedule(entry ScheduleEntry) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	heap.Push(r.scheQueue, entry)
	r.interupt()
}

func (r *ScheduleReconciler) ScheduleBatch(entries []ScheduleEntry) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, entry := range entries {
		heap.Push(r.scheQueue, entry)
	}
	r.interupt()
}

func (r *ScheduleReconciler) WaitStop() <-chan struct{} {
	c := make(chan struct{})
	defer close(c)

	return c
}
