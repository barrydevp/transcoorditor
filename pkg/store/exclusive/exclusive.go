package exclusive

import (
	// "github.com/barrydevp/lockey"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

type baseRepo struct {
	// rwLocks *lockey.RWLockKey
	rwLocks *util.RWLockKey
}

type exclusiveFunc func()

func (s *baseRepo) withLock(key string, exFn exclusiveFunc) {
	s.rwLocks.Lock(key)

	defer func() {
		s.rwLocks.Unlock(key)
	}()

	exFn()
}

func (s *baseRepo) withRLock(key string, exFn exclusiveFunc) {
	s.rwLocks.RLock(key)

	defer func() {
		s.rwLocks.RUnlock(key)
	}()

	exFn()
}

func newBaseRepo() *baseRepo {
	return &baseRepo{
		// rwLocks: lockey.NewRWLockKey(),
		rwLocks: util.NewRWLockKey(),
	}
}

// decorator pattern?
type exclusiveBackend struct {
	*store.Backend
	s store.Interface
}

func NewStore(s store.Interface) (*exclusiveBackend, error) {
	backend := &store.Backend{
		SessionImpl:     NewSession(s.Session()),
		ParticipantImpl: NewParticipant(s.Participant()),
		ReplsetImpl:     NewReplset(s.Replset()),
		LockTableImpl:   NewLockTable(s.LockTable()),
	}

	return &exclusiveBackend{
		Backend: backend,
		s:       s,
	}, nil
}

// @overide
func (s *exclusiveBackend) Close() {
	s.s.Close()
}
