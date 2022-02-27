package action

import (
	"time"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/google/uuid"
)

const (
	defaultSessionTimeout = 30
)

func NewSessionOption() *schema.SessionOptions {
	return &schema.SessionOptions{Timeout: defaultSessionTimeout}
}

func NewSession(opts *schema.SessionOptions) *schema.Session {
	now := time.Now()

	return &schema.Session{
		Id:           uuid.NewString(),
		State:        schema.SessionNew,
		Timeout:      opts.Timeout,
		CreatedAt:    &now,
		Participants: nil,
	}
}

func (ac *Action) checkSessionAvailable(s *schema.Session) error {
	if s.State == schema.SessionNew {
		return util.NewError("session was not started")
	}

	if s.IsTimeout() {
		return util.NewError("session has been expired")
	}

	if s.State != schema.SessionStarted && s.State != schema.SessionActive {
		return util.NewError("session is in %v", s.State)
	}

	return nil
}

func (ac *Action) GetSessionById(id string, populate bool) (*schema.Session, error) {
	session, err := ac.s.Session().FindById(id)

	if err != nil {
		return nil, err
	}

	// implement populate participant

	return session, nil
}

func (ac *Action) StartSession(s *schema.Session) error {
	// if s.State != schema.SessionNew {
	// 	return util.NewError("session has already started, state: %v", s.State)
	// }

	now := time.Now()

	s.State = schema.SessionStarted
	s.StartedAt = &now
	s.UpdatedAt = &now

	ac.s.Session().Save(s)

	return nil
}

func (ac *Action) JoinSession(s *schema.Session, part *schema.Participant) error {
	if err := ac.checkSessionAvailable(s); err != nil {
		return err
	}

	// omit checking participant

	part.SessionId = s.Id
	if err := ac.s.Participant().Save(part); err != nil {
		return err
	}

	if s.State == schema.SessionStarted {
		s.State = schema.SessionActive
		if err := ac.s.Session().Save(s); err != nil {
			return err
		}
	}

	return nil
}

func (ac *Action) UpdateSession(s *schema.Session, part *schema.Participant) error {
	if err := ac.checkSessionAvailable(s); err != nil {
		return err
	}

	if part.SessionId != s.Id {
		return util.NewError("this participant does not belong to updating session")
	}

	if err := ac.s.Participant().Save(part); err != nil {
		return err
	}

	// FIXME: ensuare session is in active
	if s.State == schema.SessionStarted {
		s.State = schema.SessionActive
		if err := ac.s.Session().Save(s); err != nil {
			return err
		}
	}

	return nil
}
