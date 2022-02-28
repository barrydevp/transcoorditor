package action

import (
	"time"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

func (ac *Action) checkSessionAvailable(s *schema.Session) error {
	if s.State == schema.SessionNew {
		return util.Errorf("session was not started")
	}

	if s.State != schema.SessionStarted && s.State != schema.SessionActive {
		return util.Errorf("session is in %v", s.State)
	}

	if s.IsTimeout() {
		return util.Errorf("session has been expired")
	}

	return nil
}

func (ac *Action) GetSessionById(id string, populate bool) (*schema.Session, error) {
	session, err := ac.s.Session().FindById(id)

	if err != nil {
		return nil, util.Errorf("failed to get session: %w", err)
	}

	// implement populate participant
	if populate {
		parts, err := ac.s.Participant().FindBySessionId(id)

		if err != nil {
			return nil, util.Errorf("failed to get participants: %w", err)
		}

		session.Participants = parts
	}

	return session, nil
}

func (ac *Action) StartSession(s *schema.Session) (*schema.Session, error) {
	// if s.State != schema.SessionNew {
	// 	return util.NewError("session has already started, state: %v", s.State)
	// }

	now := time.Now()

	s.State = schema.SessionStarted
	s.StartedAt = &now
	s.UpdatedAt = &now

	err := ac.s.Session().Save(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (ac *Action) JoinSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	session, err := ac.s.Session().FindById(sessionId)
	if err != nil {
		return nil, util.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, util.Errorf("session not found")
	}

	if err := ac.checkSessionAvailable(session); err != nil {
		return nil, err
	}

	if part.RequestId != "" {
		// duplicate detection by check exists participant has current requestId
		dupPart, err := ac.s.Participant().FindDupInSession(sessionId, part)

		if err != nil {
			return nil, util.Errorf("failed to check duplicate participant %w", err)
		}

		if dupPart != nil {
			// should we throw erro?
			// return nil, util.Errorf("duplicate participant %v has request %v", part.ClientId, part.RequestId)
			return dupPart, nil
		}
	}

	// part.SessionId = s.Id
	if err := ac.s.Participant().Save(part); err != nil {
		return nil, err
	}

	// first participant in session, change session State
	if session.State == schema.SessionStarted {
		session.State = schema.SessionActive
		if _, err := ac.s.Session().UpdateById(sessionId, &schema.SessionUpdate{State: &session.State}); err != nil {
			return nil, err
		}
	}

	return part, nil
}

func (ac *Action) UpdateSession(s *schema.Session, part *schema.Participant) error {
	if err := ac.checkSessionAvailable(s); err != nil {
		return err
	}

	if part.SessionId != s.Id {
		return util.Errorf("this participant does not belong to updating session")
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
