package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

// memory storage
// TBD
type Participant struct {
	m map[string]interface{}
}

func NewParticipant() *Participant {

	return &Participant{
		m: make(map[string]interface{}),
	}
}

func (s *Participant) Save(part *schema.Participant) error {
	s.m[part.Id] = part

	return nil
}

func (s *Participant) FindById(id string) (*schema.Participant, error) {
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
func (s *Participant) FindBySessionId(sessionId string) ([]*schema.Participant, error) {
	return nil, nil
}

func (s *Participant) FindDupInSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	return nil, nil
}

func (s *Participant) UpdateById(id string, partUpdate *schema.ParticipantUpdate) (*schema.Participant, error) {
	return nil, nil
}
