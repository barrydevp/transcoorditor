package action

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

var (
	ErrSessionNotFound error = errors.New("session not found")
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

func (ac *Action) findSessionById(id string) (*schema.Session, error) {
	doc, err := ac.s.Session().FindById(id)
	if err != nil {
		return nil, util.Errorf("failed to get session: %w", err)
	}

	if doc == nil {
		return nil, ErrSessionNotFound
	}

	return doc, nil
}

func (ac *Action) GetSessionById(id string, populate bool) (*schema.Session, error) {
	// @TODO: implement IO concurrent
	session, err := ac.findSessionById(id)

	if err != nil {
		return nil, err
	}

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

	if err := ac.s.Session().Save(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (ac *Action) JoinSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	session, err := ac.findSessionById(sessionId)
	if err != nil {
		return nil, err
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

	// @TODO: wrap in transaction
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

func (ac *Action) PartialCommitSession(sessionId string, partCommit *schema.ParticipantCommit) (*schema.Participant, error) {
	// @TODO: improve get session and participant concurrency
	session, err := ac.findSessionById(sessionId)
	if err != nil {
		return nil, err
	}
	if err := ac.checkSessionAvailable(session); err != nil {
		return nil, err
	}

	part, err := ac.findParticipantById(*partCommit.Id)
	if err != nil {
		return nil, err
	}

	if part.SessionId != session.Id {
		return nil, util.Errorf("commit participant not belong to this session")
	}
	if part.State != schema.ParticipantActive {
		return nil, util.Errorf("this participant has already committed, current state: %v", part.State)
	}

	part.State = schema.ParticipantCommitted

	if partCommit.Compensate != nil {
		partCommit.Compensate.Status = schema.PartActionPending
		partCommit.Compensate.InvokedCount = 0
	}
	if partCommit.Complete != nil {
		partCommit.Complete.Status = schema.PartActionPending
		partCommit.Complete.InvokedCount = 0
	}

	partUpdate := &schema.ParticipantUpdate{
		State:            &part.State,
		CompensateAction: partCommit.Compensate,
		CompleteAction:   partCommit.Complete,
	}

	// @TODO: wrap in transaction
	part, err = ac.s.Participant().UpdateById(*partCommit.Id, partUpdate)
	if err != nil {
		return nil, util.Errorf("failed to commit participant: %w", err)
	}

	return part, nil
}

func (ac *Action) CommitSession(id string) (*schema.Session, error) {
	session, err := ac.GetSessionById(id, true)
	if err != nil {
		return nil, err
	}

	if err = session.AbleToCommitOrRollback(); err != nil {
		return nil, err
	}

	session.State = schema.SessionCommitting
	if _, err := ac.s.Session().UpdateById(id, &schema.SessionUpdate{State: &session.State}); err != nil {
		return nil, err
	}

	errs := ac.handlePartComplete(session)
	if len(errs) > 0 {
		session.State = schema.SessionCommitFailed
		session.Errors = errs
	} else {
		session.State = schema.SessionCommitted
	}

	if _, err := ac.s.Session().UpdateById(id, &schema.SessionUpdate{State: &session.State, Errors: &session.Errors}); err != nil {
		return nil, err
	}

	return session, err
}

func (ac *Action) abortSession(session *schema.Session) (*schema.Session, error) {
	if err := session.AbleToCommitOrRollback(); err != nil {
		return nil, err
	}

	session.State = schema.SessionAborting
	if _, err := ac.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State}); err != nil {
		return nil, err
	}

	errs := ac.handlePartCompensate(session)
	if len(errs) > 0 {
		session.State = schema.SessionAbortFailed
		session.Errors = errs
	} else {
		session.State = schema.SessionAborted
	}

	if _, err := ac.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State, Errors: &session.Errors}); err != nil {
		return nil, err
	}

	return session, nil
}

func (ac *Action) AbortSession(id string) (*schema.Session, error) {
	session, err := ac.GetSessionById(id, true)
	if err != nil {
		return nil, err
	}

	return ac.abortSession(session)

}
