package mongodb

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

// memory storage
// TBD
type Participant struct {
	m map[string]interface{}
	o *StoreOptions
}

func NewParticipant(opts *StoreOptions) *Participant {

	return &Participant{
		m: make(map[string]interface{}),
		o: opts,
	}
}

func (s *Participant) Save(part *schema.Participant) error {
	s.m[part.Id] = part

	return nil
}

func (s *Participant) FindById(id string) (*schema.Participant, error) {
	data := s.m[id]

	if data == nil {
		return nil, util.NewError("not found")
	}

	part, ok := data.(*schema.Participant)
	if !ok {
		delete(s.m, id)

		return nil, util.NewError("detect unexpected behavior")
	}

	return part, nil
}
