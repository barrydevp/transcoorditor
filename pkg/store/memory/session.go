package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

// memory storage
// TBD
type sessionRepo struct {
	m map[string]interface{}
}

func NewSession() *sessionRepo {

	return &sessionRepo{
		m: make(map[string]interface{}),
	}
}

func (s *sessionRepo) Save(session *schema.Session) error {
	s.m[session.Id] = session

	return nil
}

func (s *sessionRepo) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	return nil, nil
}

func (s *sessionRepo) Find(search *schema.SessionSearch) ([]*schema.Session, error) {
	return nil, nil
}

func (s *sessionRepo) FindAllUnfinished() ([]*schema.Session, error) {
	return nil, nil
}

func (s *sessionRepo) FindById(id string) (*schema.Session, error) {
	data := s.m[id]

	if data == nil {
		return nil, util.Errorf("not found")
	}

	session, ok := data.(*schema.Session)
	if !ok {
		delete(s.m, id)

		return nil, util.Errorf("detect unexpected behavior")
	}

	return session, nil
}

func (s *sessionRepo) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (*schema.Session, error) {
	return nil, nil
}

func (s *sessionRepo) DeleteById(id string) (*schema.Session, error) {
	return nil, nil
}
