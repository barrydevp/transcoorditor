package schema

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/google/uuid"
)

type SessionState string

const (
	SessionNew             SessionState = "New"
	SessionStarted         SessionState = "Started"
	SessionActive          SessionState = "Active"
	SessionCommitting      SessionState = "Committing"
	SessionCommitted       SessionState = "Commited"
	SessionCommitFailed    SessionState = "CommitFailed"
	SessionAborting        SessionState = "Aborting"
	SessionAborted         SessionState = "Aborted"
	SessionAbortFailed     SessionState = "AbortFailed"
	SessionTerminating     SessionState = "Terminating"
	SessionTerminated      SessionState = "Terminated" // timeout session was auto terminated
	SessionTerminateFailed SessionState = "TerminateFailed"

	SessionMaximumRetry = 5
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
	Retries   int          `json:"retries" bons:"retries"`

	// for edges field (relations associate field)
	Participants []*Participant `json:"participants,omitempty" bson:"-"`
}

type SessionUpdate struct {
	State     *SessionState
	Errors    *[]string
	Timeout   *int `json:"timeout"`
	UpdatedAt *time.Time
	StartedAt *time.Time
	Retries   *int
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

func (s *Session) TimedoutAt() time.Time {
	return s.StartedAt.Add(time.Second * time.Duration(s.Timeout))
}

func (s *Session) IsTimeout() bool {
	if s.StartedAt == nil {
		return false
	}

	return time.Now().After(s.TimedoutAt())
}

func (s *Session) IsFinished() bool {
	switch s.State {
	case SessionCommitted, SessionAborted, SessionTerminated:
		return true
	default:
		return false
	}
}

func (s *Session) IsMaximumRetry() bool {
	return s.Retries > SessionMaximumRetry
}

func UnfinishedSessionStates() []string {
	return []string{
		string(SessionStarted),
		string(SessionActive),
		string(SessionCommitting),
		string(SessionCommitFailed),
		string(SessionAborting),
		string(SessionAbortFailed),
		string(SessionTerminating),
		string(SessionTerminateFailed),
	}
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

func (s *Session) CheckAllPartAbleToEnd() error {
	for _, part := range s.Participants {

		switch part.State {
		case ParticipantActive:
			return util.Errorf("some participant is in Active (in-processing and hasn't commit)")
		case ParticipantCompensating:
			return util.Errorf("some participant is in Compensating")
		case ParticipantCompleting:
			return util.Errorf("some participant is in Completing")
		}
	}

	return nil
}

var (
	ErrSessionIsCommitting    = errors.New("session is Committing")
	ErrSessionWasCommitted    = errors.New("session was Committed")
	ErrSessionWasCommitFailed = errors.New("session was CommitFailed")
	ErrSessionIsAborting      = errors.New("session is Aborting")
	ErrSessionWasAborted      = errors.New("session was Aborted")
	ErrSessionWasAbortFailed  = errors.New("session was AbortFailed")
)

func (s *Session) AbleToCommitOrRollback() error {
	switch s.State {
	case SessionCommitting:
		return ErrSessionIsCommitting
	case SessionCommitted:
		return ErrSessionWasCommitted
	case SessionCommitFailed:
		return ErrSessionWasCommitFailed
	case SessionAborting:
		return ErrSessionIsAborting
	case SessionAborted:
		return ErrSessionWasAborted
	case SessionAbortFailed:
		return ErrSessionWasAbortFailed

	}

	if err := s.CheckSessionActive(); err != nil {
		return err
	}

	if err := s.CheckAllPartAbleToEnd(); err != nil {
		return err
	}

	return nil
}

func (s *Session) GetParticipantAt(id int64) *Participant {
	if len(s.Participants) < int(id) || s.Participants[id-1] == nil {
		return nil
	}

	return s.Participants[id-1]
}

type SessionSearch struct {
	State *string
}
