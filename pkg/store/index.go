package store

import "github.com/barrydevp/transcoorditor/pkg/schema"

type Interface interface {
	Session() Session
	Participant() Participant
    Close()
}

type Session interface {
	Save(s *schema.Session) error
	FindById(id string) (*schema.Session, error)
}

type Participant interface {
	Save(s *schema.Participant) error
	FindById(id string) (*schema.Participant, error)
}

