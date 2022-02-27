package store

import (
	"github.com/barrydevp/transcoorditor/pkg/store/memory"
)

type Memory struct {
	session     *memory.Session
	participant *memory.Participant
}

func NewMemoryStore() (Interface, error) {

	return &Memory{
		session:     memory.NewSession(),
		participant: memory.NewParticipant(),
	}, nil
}

func (s *Memory) Close() {

}

func (s *Memory) Session() Session {
	return s.session
}

func (s *Memory) Participant() Participant {
	return s.participant
}
