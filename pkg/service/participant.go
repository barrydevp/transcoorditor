package service

import (
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
)

var (
	ErrParticipantNotFound = exception.AppNotFoundf("participant not found")
)

func (srv *Service) findParticipantById(sessionId string, id int64) (*schema.Participant, error) {
	doc, err := srv.s.Participant().FindBySessionAndId(sessionId, id)
	if err != nil {
		return nil, exception.Errorf("failed to get participant: %w", err)
	}

	if doc == nil {
		return nil, ErrParticipantNotFound
	}

	return doc, nil
}

type PartActionHandler func(*schema.Participant) (*schema.ParticipantUpdate, error)

func (srv *Service) handlePartAction(session *schema.Session, handler PartActionHandler) []string {
	if len(session.Participants) == 0 {
		return nil
	}

	var errs []string

	for _, part := range session.Participants {
		partUpdate, err := handler(part)

		if err != nil {
			errs = append(errs, err.Error())
		}

		// srv.l.Info(partUpdate)

		if partUpdate != nil {
			if _, err = srv.s.Participant().UpdateBySessionAndId(session.Id, part.Id, partUpdate); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	return errs
}

func (srv *Service) handlePartComplete(session *schema.Session) []string {
	return srv.handlePartAction(session, func(part *schema.Participant) (*schema.ParticipantUpdate, error) {
		var err error

		completeAc := part.CompleteAction
		partState := schema.ParticipantCompleted

		if completeAc != nil {

			// if completeAc.IsFinished() {
			// 	return nil, nil
			// }

			err = completeAc.InvokePartAction()
			if err != nil {
				partState = schema.ParticipantCompleteFailed
			}
		}

		return &schema.ParticipantUpdate{
			State:          &partState,
			CompleteAction: completeAc,
		}, err
	})
}

func (srv *Service) handlePartCompensate(session *schema.Session) []string {
	return srv.handlePartAction(session, func(part *schema.Participant) (*schema.ParticipantUpdate, error) {
		var err error

		compensateAc := part.CompensateAction
		partState := schema.ParticipantCompensated

		if compensateAc != nil {

			// if compensateAction.IsFinished() {
			// 	continue
			// }

			err = compensateAc.InvokePartAction()
			if err != nil {
				partState = schema.ParticipantCompensateFailed
			}
		}

		return &schema.ParticipantUpdate{
			State:            &partState,
			CompensateAction: compensateAc,
		}, err
	})
}
