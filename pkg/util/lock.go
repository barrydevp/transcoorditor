package util

import (
	// "fmt"
	"sync"
	"sync/atomic"
)

type rwLock struct {
	*sync.RWMutex
	ownerCount int32
}

type RWLockKey struct {
	rw    sync.RWMutex
	locks map[string]*rwLock
}

func NewRWLockKey() *RWLockKey {
	return &RWLockKey{
		rw:    sync.RWMutex{},
		locks: make(map[string]*rwLock),
	}
}

// get the lock for Lock
func (l *RWLockKey) getLock(key string) *rwLock {
	// read phase, retrive lock if exists
	l.rw.RLock()
	lock, ok := l.locks[key]

	if ok {
		atomic.AddInt32(&lock.ownerCount, 1)
		l.rw.RUnlock()

		return lock
	}
	l.rw.RUnlock() // unlock read phase

	// write phase, create new lock
	l.rw.Lock()

	// re-check if multiple go-routine enter write phase and one onother has already created new lock
	lock, ok = l.locks[key]
	if ok {
		atomic.AddInt32(&lock.ownerCount, 1)

		l.rw.Unlock()
		return lock
	}

	// actually create new lock
	lock = &rwLock{
		RWMutex:    &sync.RWMutex{},
		ownerCount: 1,
	}
	l.locks[key] = lock

	l.rw.Unlock()

	return lock
}

// get lock for Unlock and put the lock back to pool, if lock's ownerCount == 0, mean no other has owned the lock, so we will remove it
func (l *RWLockKey) putLock(key string) *rwLock {
	l.rw.Lock()

	lock, ok := l.locks[key]

	// unexpected behavior, the lock has been released by another goroutine or unlock twice
	if !ok {
		return nil
	}

	// check if lock's ownerCount == 0, no other has using the lock, remove it from pool
	lock.ownerCount--
	if lock.ownerCount == 0 {
		delete(l.locks, key)
	}

	l.rw.Unlock()

	return lock
}

func (l *RWLockKey) Lock(key string) {
	// fmt.Println("lock0", key, l.locks)
	l.getLock(key).Lock()
	// fmt.Println("lock", key,  l.locks)
}

func (l *RWLockKey) Unlock(key string) {
	// fmt.Println("unlock0", key, l.locks)
	if lock := l.putLock(key); lock != nil {
		lock.Unlock()
	}
	// fmt.Println("unlock", key, l.locks)
}

func (l *RWLockKey) RLock(key string) {
	l.getLock(key).RLock()
}

func (l *RWLockKey) RUnlock(key string) {
	if lock := l.putLock(key); lock != nil {
		lock.RUnlock()
	}
}
