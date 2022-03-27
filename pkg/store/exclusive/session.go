package exclusive

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type sessionRepo struct {
	*baseRepo
	s store.Session
}

func NewSession(s store.Session) store.Session {
	return &sessionRepo{
		baseRepo: newBaseRepo(),
		s:        s,
	}
}

func (s *sessionRepo) Save(session *schema.Session) error {
	return s.s.Save(session)
}

func (s *sessionRepo) PutById(id string, schemaUpdate *schema.Session) (session *schema.Session, err error) {
	s.withLock(id, func() {
		session, err = s.s.PutById(id, schemaUpdate)
	})

	return
}

func (s *sessionRepo) Find(search *schema.SessionSearch) (sessions []*schema.Session, err error) {
	return s.s.Find(search)
}

func (s *sessionRepo) FindAllUnfinished() ([]*schema.Session, error) {
	return s.s.FindAllUnfinished()
}

func (s *sessionRepo) FindById(id string) (session *schema.Session, err error) {
	s.withLock(id, func() {
		session, err = s.s.FindById(id)
	})

	return
}

func (s *sessionRepo) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (session *schema.Session, err error) {
	s.withLock(id, func() {
		session, err = s.s.UpdateById(id, schemaUpdate)
	})

	return
}

func (s *sessionRepo) DeleteById(id string) (session *schema.Session, err error) {
	s.withLock(id, func() {
		session, err = s.s.DeleteById(id)

	})

	return
}
