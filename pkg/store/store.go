package store

import "github.com/barrydevp/transcoorditor/pkg/schema"

type Interface interface {
	Session() Session
	Participant() Participant
	Close()
}

type Session interface {
	Save(s *schema.Session) error
	PutById(id string, session *schema.Session) (*schema.Session, error)
	FindById(id string) (*schema.Session, error)
	UpdateById(id string, update *schema.SessionUpdate) (*schema.Session, error)
}

type Participant interface {
	Save(part *schema.Participant) error
	PutBySessionAndId(sessionId string, id int64, part *schema.Participant) (*schema.Participant, error)
	FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error)
	FindBySessionId(sessionId string) ([]*schema.Participant, error)
	FindDupInSession(sesionId string, part *schema.Participant) (*schema.Participant, error)
	UpdateBySessionAndId(sessionId string, id int64, update *schema.ParticipantUpdate) (*schema.Participant, error)
	CountBySessionId(sessionId string) (int64, error)
}
