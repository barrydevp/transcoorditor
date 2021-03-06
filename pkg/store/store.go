package store

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/hashicorp/raft"
)

var (
	ErrSessionNotFound = exception.AppNotFoundf("session not found")
	ErrLockNotFound    = errors.New("lock not found")
	ErrLockNotOwner    = errors.New("lock was belong to another onwer")
	ErrLockExpired     = errors.New("lock has been expired")
	ErrLockExists      = errors.New("lock has been exist and not expired yet")
)

type (
	Interface interface {
		Session() Session
		Participant() Participant
		Replset() Replset
		LockTable() LockTable
		GetApplier() cluster.Applier
		Close()
	}

	Replset interface {
		SaveLastLog(log *raft.Log) error
		GetLastLog() (*raft.Log, error)
	}

	Session interface {
		Save(s *schema.Session) error
		PutById(id string, session *schema.Session) (*schema.Session, error)
		FindById(id string) (*schema.Session, error)
		Find(search *schema.SessionSearch) ([]*schema.Session, error)
		FindAllUnfinished() ([]*schema.Session, error)
		UpdateById(id string, update *schema.SessionUpdate) (*schema.Session, error)
		DeleteById(id string) (*schema.Session, error)
	}

	Participant interface {
		Save(part *schema.Participant) error
		PutBySessionAndId(sessionId string, id int64, part *schema.Participant) (*schema.Participant, error)
		FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error)
		FindBySessionId(sessionId string) ([]*schema.Participant, error)
		FindDupInSession(sesionId string, part *schema.Participant) (*schema.Participant, error)
		UpdateBySessionAndId(sessionId string, id int64, update *schema.ParticipantUpdate) (*schema.Participant, error)
		CountBySessionId(sessionId string) (int64, error)
		DeleteBySessionId(sessionId string) (int64, error)
	}

	LockTable interface {
		Save(l *schema.LockEntry) error
		Update(l *schema.LockEntry) error
		Find(key string) (*schema.LockEntry, error)
		FindWithOwner(key string, owner string) (*schema.LockEntry, error)
		Delete(l *schema.LockEntry) error
		DeleteByOwner(owner string) (int64, error)
	}
)

type Backend struct {
	SessionImpl     Session
	ParticipantImpl Participant
	ReplsetImpl     Replset
	LockTableImpl   LockTable
}

func (b *Backend) Session() Session {
	return b.SessionImpl
}

func (b *Backend) Participant() Participant {
	return b.ParticipantImpl
}

func (b *Backend) Replset() Replset {
	return b.ReplsetImpl
}

func (b *Backend) LockTable() LockTable {
	return b.LockTableImpl
}

func (b *Backend) GetApplier() cluster.Applier {
	return nil
}
