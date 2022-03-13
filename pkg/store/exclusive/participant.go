package exclusive

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type participantRepo struct {
	*baseRepo
	s store.Participant
}

func NewParticipant(participant store.Participant) *participantRepo {
	return &participantRepo{
		baseRepo: newBaseRepo(),
		s:        participant,
	}
}

func (s *participantRepo) Save(part *schema.Participant) error {
	return s.s.Save(part)
}

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, part *schema.Participant) (pa *schema.Participant, err error) {
	s.withLock(sessionId, func() {
		pa, err = s.s.PutBySessionAndId(sessionId, id, part)
	})

	return
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (part *schema.Participant, err error) {
	s.withLock(sessionId, func() {
		part, err = s.s.FindBySessionAndId(sessionId, id)
	})

	return
}

func (s *participantRepo) FindBySessionId(sessionId string) (parts []*schema.Participant, err error) {
	s.withLock(sessionId, func() {
		parts, err = s.s.FindBySessionId(sessionId)
	})

	return
}

func (s *participantRepo) FindDupInSession(sessionId string, part *schema.Participant) (pa *schema.Participant, err error) {
	s.withLock(sessionId, func() {
		pa, err = s.s.FindDupInSession(sessionId, part)
	})

	return
}

func (s *participantRepo) UpdateBySessionAndId(sessionId string, id int64, partUpdate *schema.ParticipantUpdate) (part *schema.Participant, err error) {
	s.withLock(sessionId, func() {
		part, err = s.s.UpdateBySessionAndId(sessionId, id, partUpdate)
	})

	return
}

func (s *participantRepo) CountBySessionId(sessionId string) (count int64, err error) {
	s.withLock(sessionId, func() {
		count, err = s.s.CountBySessionId(sessionId)
	})

	return
}

func (s *participantRepo) DeleteBySessionId(sessionId string) (count int64, err error) {
	s.withLock(sessionId, func() {
		count, err = s.s.DeleteBySessionId(sessionId)
	})

	return
}
