package schema

import (
	"time"
)

type SessionState string

const (
	SessionNew         SessionState = "New"
	SessionStarted                  = "Started"
	SessionActive                   = "Active"
	SessionCommitting               = "Committing"
	SessionCommitted                = "Commited"
	SessionAborting                 = "Aborting"
	SessionAborted                  = "Aborted"
	SessionTerminating              = "Terminating"
	SessionTerminated               = "Terminated"
)

type SessionOptions struct {
	Timeout int `json:"timeout"`
}

type Session struct {
	// represents storage field. eg: mongodb field, mysql column
	Id string `json:"id" bson:"id"`

	State     SessionState `json:"state" bson:"state"`
	Timeout   int          `json:"timeout" bson:"timeout"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	StartedAt *time.Time   `json:"startedAt,omitempty" bson:"startedAt,omitempty"`
	CreatedAt *time.Time   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`

	// for edges field (relations associate field)
	Participants []*Participant `json:"participants,omitempty" bson:"-"`
}

func (s *Session) IsTimeout() bool {
	if s.StartedAt == nil {
		return false
	}

	return time.Now().After(s.StartedAt.Add(time.Second * time.Duration(s.Timeout)))
}
