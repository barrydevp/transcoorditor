package schema

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/google/uuid"
)

type SessionState string

const (
	SessionNew          SessionState = "New"
	SessionStarted                   = "Started"
	SessionActive                    = "Active"
	SessionCommitting                = "Committing"
	SessionCommitted                 = "Commited"
	SessionCommitFailed              = "CommitFailed"
	SessionAborting                  = "Aborting"
	SessionAborted                   = "Aborted"
	SessionAbortFailed               = "AbortFailed"
)

type SessionOptions struct {
	Timeout int `json:"timeout"`
}

const (
	defaultSessionTimeout = 120 // 2 mins
)

var (
	ErrSessionExpired = errors.New("session has been expired")
)

func NewSessionOption() *SessionOptions {
	return &SessionOptions{Timeout: defaultSessionTimeout}
}

type Session struct {
	// represents storage field. eg: mongodb field, mysql column
	Id string `json:"id" bson:"id"`

	State     SessionState `json:"state" bson:"state"`
	Timeout   int          `json:"timeout" bson:"timeout"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	StartedAt *time.Time   `json:"startedAt,omitempty" bson:"startedAt,omitempty"`
	CreatedAt *time.Time   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Errors    []string     `json:"errors,omitempty" bson:"errors,omitempty"`

	// for edges field (relations associate field)
	Participants []*Participant `json:"participants,omitempty" bson:"-"`
}

type SessionUpdate struct {
	State     *SessionState
	Errors    *[]string
	Timeout   *int `json:"timeout"`
	UpdatedAt *time.Time
	StartedAt *time.Time
}

func NewSession(opts *SessionOptions) *Session {
	now := time.Now()

	return &Session{
		Id:           uuid.NewString(),
		State:        SessionNew,
		Timeout:      opts.Timeout,
		CreatedAt:    &now,
		Participants: nil,
	}
}

func (s *Session) IsTimeout() bool {
	if s.StartedAt == nil {
		return false
	}

	return time.Now().After(s.StartedAt.Add(time.Second * time.Duration(s.Timeout)))
}

func (s *Session) CheckSessionActive() error {
	if s.State == SessionNew {
		return util.Errorf("session was not started")
	}

	if s.State != SessionStarted && s.State != SessionActive {
		return util.Errorf("session is not Active, current is %v", s.State)
	}

	if s.IsTimeout() {
		return ErrSessionExpired
	}

	return nil
}

func (s *Session) IsAllPartAbleToEnd() bool {
	for _, part := range s.Participants {

		switch part.State {
		case ParticipantActive, ParticipantCompensating, ParticipantCompleting:
			return false
		}
	}

	return true
}

func (s *Session) AbleToCommitOrRollback() error {
	// if s.State != SessionActive && s.State != SessionStarted {
	// 	return util.Errorf("session is not Active, current is %v", s.State)
	// }
	//
	// if s.IsTimeout() {
	// 	return ErrSessionExpired
	// }

	if err := s.CheckSessionActive(); err != nil {
		return err
	}

	if !s.IsAllPartAbleToEnd() {
		return util.Errorf("session has some participant not able to end")
	}

	return nil
}
