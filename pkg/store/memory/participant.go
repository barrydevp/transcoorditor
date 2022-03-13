package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

// memory storage
// TBD
type participantRepo struct {
	m map[int64]interface{}
}

func NewParticipant() *participantRepo {

	return &participantRepo{
		m: make(map[int64]interface{}),
	}
}

func (s *participantRepo) Save(part *schema.Participant) error {
	s.m[part.Id] = part

	return nil
}

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, part *schema.Participant) (*schema.Participant, error) {
	return nil, nil
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error) {
	data := s.m[id]

	if data == nil {
		return nil, util.Errorf("not found")
	}

	part, ok := data.(*schema.Participant)
	if !ok {
		delete(s.m, id)

		return nil, util.Errorf("detect unexpected behavior")
	}

	return part, nil
}

// @TODO
func (s *participantRepo) FindBySessionId(sessionId string) ([]*schema.Participant, error) {
	return nil, nil
}

func (s *participantRepo) FindDupInSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	return nil, nil
}

func (s *participantRepo) UpdateBySessionAndId(sessionId string, id int64, partUpdate *schema.ParticipantUpdate) (*schema.Participant, error) {
	return nil, nil
}

func (s *participantRepo) CountBySessionId(sessionId string) (int64, error) {
	return 0, nil
}

func (s *participantRepo) DeleteBySessionId(sessionId string) (int64, error) {
	return 0, nil
}
