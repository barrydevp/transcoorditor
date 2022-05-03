package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type memoryBackend struct {
	*store.Backend
}

func NewStore() (store.Interface, error) {
	backend := &store.Backend{
		SessionImpl:     NewSession(),
		ParticipantImpl: NewParticipant(),
		LockTableImpl: NewLockTable(),
	}

	return &memoryBackend{
		Backend: backend,
	}, nil
}

func (s *memoryBackend) Close() {

}
