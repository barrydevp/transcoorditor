package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/sirupsen/logrus"
)

// memory storage
// TBD
type Session struct {
	m map[string]interface{}
	l *logrus.Entry
}

func NewSession() *Session {

	return &Session{
		m: make(map[string]interface{}),
		l: common.Logger().WithFields(logrus.Fields{
			"store": "memory",
			"name":  "session",
		}),
	}
}

func (s *Session) Save(session *schema.Session) error {
	s.m[session.Id] = session

	return nil
}

func (s *Session) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	return nil, nil
}

func (s *Session) Find(search *schema.SessionSearch) ([]*schema.Session, error) {
	return nil, nil
}

func (s *Session) FindById(id string) (*schema.Session, error) {
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

func (s *Session) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (*schema.Session, error) {
	return nil, nil
}

func (s *Session) DeleteById(id string) (*schema.Session, error) {
	return nil, nil
}
