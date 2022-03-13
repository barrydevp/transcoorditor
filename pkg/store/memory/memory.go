package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type memoryStore struct {
	session     *sessionRepo
	participant *participantRepo
}

func NewStore() (store.Interface, error) {

	return &memoryStore{
		session:     NewSession(),
		participant: NewParticipant(),
	}, nil
}

func (s *memoryStore) Close() {

}

func (s *memoryStore) Session() store.Session {
	return s.session
}

func (s *memoryStore) Participant() store.Participant {
	return s.participant
}
