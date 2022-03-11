package service

import (
	"time"

	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
)

var (
	ErrSessionNotFound error = exception.Errorf("%w: session was not found in storage", exception.ErrNotFound)
)

func (srv *Service) findSessionById(id string) (*schema.Session, error) {
	doc, err := srv.s.Session().FindById(id)
	if err != nil {
		return nil, exception.Errorf("failed to get session: %w", err)
	}

	return doc, nil
}

func (srv *Service) GetSessionById(id string, populate bool) (*schema.Session, error) {
	// @TODO: implement IO concurrent
	session, err := srv.findSessionById(id)

	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, ErrSessionNotFound
	}

	if populate {
		parts, err := srv.s.Participant().FindBySessionId(id)

		if err != nil {
			return nil, exception.Errorf("failed to get participants: %w", err)
		}

		session.Participants = parts
	}

	return session, nil
}

func (srv *Service) ListSession() ([]*schema.Session, error) {
	docs, err := srv.s.Session().Find(&schema.SessionSearch{})

	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (srv *Service) DeleteSessionById(id string) (*schema.Session, error) {
	session, err := srv.s.Session().DeleteById(id)
	if err != nil {
		return nil, exception.Errorf("failed to delete session: %w", err)
	}

	_, err = srv.s.Participant().DeleteBySessionId(id)
	if err != nil {
		return nil, exception.Errorf("failed to delete participants of session: %w", err)
	}

	srv.l.Info("participant deleted count: ", 10)

	if session == nil {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

func (srv *Service) PutSessionById(s *schema.Session) (*schema.Session, error) {
	session, err := srv.s.Session().PutById(s.Id, s)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (srv *Service) StartSession(s *schema.Session) (*schema.Session, error) {
	// if s.State != schema.SessionNew {
	// 	return util.NewError("session has already started, state: %v", s.State)
	// }

	now := time.Now()

	s.State = schema.SessionStarted
	s.StartedAt = &now
	s.UpdatedAt = &now

	if err := srv.s.Session().Save(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (srv *Service) JoinSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	session, err := srv.findSessionById(sessionId)
	if err != nil {
		return nil, err
	}

	if err := session.CheckSessionActive(); err != nil {
		return nil, exception.Errorw(err, exception.ErrPreconditionFailed)
	}

	if part.RequestId != "" {
		// duplicate detection by check exists participant has current requestId
		dupPart, err := srv.s.Participant().FindDupInSession(sessionId, part)
		if err != nil {
			return nil, exception.Errorf("failed to check duplicate participant %w", err)
		}

		if dupPart != nil {
			// should we throw erro?
			// return nil, exception.Errorf("duplicate participant %v has request %v", part.ClientId, part.RequestId)
			return dupPart, nil
		}
	}

	partNum, err := srv.s.Participant().CountBySessionId(sessionId)
	if err != nil {
		return nil, exception.Errorf("failed to get number participant in session %w", err)
	}

	// @TODO: wrap in transaction
	// part.SessionId = s.Id
	part.Id = partNum + 1
	if err := srv.s.Participant().Save(part); err != nil {
		return nil, err
	}

	// first participant in session, change session State
	if session.State == schema.SessionStarted {
		session.State = schema.SessionActive
		if _, err := srv.s.Session().UpdateById(sessionId, &schema.SessionUpdate{State: &session.State}); err != nil {
			return nil, err
		}
	}

	return part, nil
}

func (srv *Service) PartialCommitSession(sessionId string, partCommit *schema.ParticipantCommit) (*schema.Participant, error) {
	// @TODO: improve get session and participant concurrency
	session, err := srv.findSessionById(sessionId)
	if err != nil {
		return nil, err
	}
	if err := session.CheckSessionActive(); err != nil {
		return nil, err
	}

	part, err := srv.findParticipantById(sessionId, *partCommit.Id)
	if err != nil {
		return nil, err
	}

	if part.SessionId != session.Id {
		return nil, exception.Errorf("commit participant not belong to this session. %w", exception.ErrPreconditionFailed)
	}

	if part.State != schema.ParticipantActive {
		return nil, exception.Errorf("this participant has already committed, current state: %v. %w", part.State, exception.ErrPreconditionFailed)
	}

	part.State = schema.ParticipantCommitted

	if partCommit.Compensate != nil {
		partCommit.Compensate.Status = schema.PartActionCreated
		partCommit.Compensate.InvokedCount = 0
	}
	if partCommit.Complete != nil {
		partCommit.Complete.Status = schema.PartActionCreated
		partCommit.Complete.InvokedCount = 0
	}

	partUpdate := &schema.ParticipantUpdate{
		State:            &part.State,
		CompensateAction: partCommit.Compensate,
		CompleteAction:   partCommit.Complete,
	}

	// @TODO: wrap in transaction
	part, err = srv.s.Participant().UpdateBySessionAndId(sessionId, *partCommit.Id, partUpdate)
	if err != nil {
		return nil, exception.Errorf("failed to commit participant: %w", err)
	}

	return part, nil
}

func (srv *Service) commitSession(session *schema.Session) (*schema.Session, error) {
	var err error = nil

	if err = session.AbleToCommitOrRollback(); err != nil {
		return nil, exception.Errorw(err, exception.ErrPreconditionFailed)
	}

	session.State = schema.SessionCommitting
	if _, err = srv.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State}); err != nil {
		return nil, err
	}

	errs := srv.handlePartComplete(session)
	if len(errs) > 0 {
		session.State = schema.SessionCommitFailed
		session.Errors = errs
		err = exception.Errorf("failed to handle CompleteAction on participants")
	} else {
		session.State = schema.SessionCommitted
	}

	if _, err = srv.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State, Errors: &session.Errors}); err != nil {
		return nil, err
	}

	return session, err
}

func (srv *Service) CommitSession(id string) (*schema.Session, error) {
	session, err := srv.GetSessionById(id, true)
	if err != nil {
		return nil, err
	}

	return srv.commitSession(session)
}

func (srv *Service) abortSession(session *schema.Session) (*schema.Session, error) {
	var err error = nil

	if err = session.AbleToCommitOrRollback(); err != nil {
		return nil, exception.Errorw(err, exception.ErrPreconditionFailed)
	}

	session.State = schema.SessionAborting
	if _, err := srv.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State}); err != nil {
		return nil, err
	}

	errs := srv.handlePartCompensate(session)
	if len(errs) > 0 {
		session.State = schema.SessionAbortFailed
		session.Errors = errs
		err = exception.Errorf("failed to handle CompensateAction on participants")
	} else {
		session.State = schema.SessionAborted
	}

	if _, err := srv.s.Session().UpdateById(session.Id, &schema.SessionUpdate{State: &session.State, Errors: &session.Errors}); err != nil {
		return nil, err
	}

	return session, err
}

func (srv *Service) AbortSession(id string) (*schema.Session, error) {
	session, err := srv.GetSessionById(id, true)
	if err != nil {
		return nil, err
	}

	return srv.abortSession(session)
}
