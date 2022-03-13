package exclusive

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

type baseRepo struct {
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
		rwLocks: util.NewRWLockKey(),
	}
}

// decorator pattern?
type exclusiveStore struct {
	s store.Interface

	session     *sessionRepo
	participant *participantRepo
}

func NewStore(store store.Interface) (store.Interface, error) {
	return &exclusiveStore{
		session:     NewSession(store.Session()),
		participant: NewParticipant(store.Participant()),
	}, nil
}

// @overide
func (s *exclusiveStore) Close() {

}

func (s *exclusiveStore) Session() store.Session {
	return s.session
}

func (s *exclusiveStore) Participant() store.Participant {
	return s.participant
}
