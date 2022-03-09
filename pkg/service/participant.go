package service

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrParticipantNotFound = errors.New("participant not found")
	ErrInvalidActionUri    = errors.New("invalid action's uri")
	ErrActionRequestFailed = errors.New("action request failed")
)

func (srv *Service) findParticipantById(sessionId string, id int64) (*schema.Participant, error) {
	doc, err := srv.s.Participant().FindBySessionAndId(sessionId, id)
	if err != nil {
		return nil, util.Errorf("failed to get participant: %w", err)
	}

	if doc == nil {
		return nil, ErrParticipantNotFound
	}

	return doc, nil
}

func (srv *Service) invokePartAction(action *schema.ParticipantAction) (err error) {
	var (
		req  *resty.Request
		resp *resty.Response
	)

	action.InvokedCount++

	result := &schema.PartActionResult{}

	if action.Uri == nil {
		err = ErrInvalidActionUri
		goto RET
	}

	// build request
	req = util.GetRequest().R()

	if action.Data != nil {
		// action.Data is interface{} so in case of data retrieve from mongodb, Data may be bson.D
		// bson.D is slice so it convert into Array, so we need convert into bson.M for json Object
		body := action.Data
		switch v := action.Data.(type) {
		case primitive.D:
			body = v.Map()
		}

		req.SetBody(body)
	}

	resp, err = req.Post(*action.Uri)

	if err != nil {
		goto RET
	}

	result.StatusCode = resp.StatusCode()
	result.Status = resp.Status()
	result.Proto = resp.Proto()
	result.Time = resp.Time()
	result.ReceivedAt = resp.ReceivedAt()
	result.Body = resp.String()

	if resp.StatusCode() != 200 {
		err = ErrActionRequestFailed
		goto RET
	}

RET:
	if err != nil {
		result.Error = err.Error()
	}

	action.Results = append(action.Results, result)

	return err
}

func (srv *Service) handlePartComplete(session *schema.Session) []string {
	if len(session.Participants) == 0 {
		return nil
	}

	var errs []string

	for _, part := range session.Participants {
		partUpdate := &schema.ParticipantUpdate{
			State:          &part.State,
			CompleteAction: part.CompleteAction,
		}

		if part.CompleteAction != nil {
			completeAc := part.CompleteAction

			if completeAc.IsFinished() {
				continue
			}

			// @TODO: validate action

			// @TODO: update to processing

			err := srv.invokePartAction(completeAc)
			if err != nil {
				completeAc.Status = schema.PartActionFailed
				part.State = schema.ParticipantCompleteFailed
				errs = append(errs, err.Error())
			} else {
				completeAc.Status = schema.PartActionCompleted
				part.State = schema.ParticipantCompleted
			}
		} else {
			part.State = schema.ParticipantCompleted
		}

		if _, err := srv.s.Participant().UpdateBySessionAndId(session.Id, part.Id, partUpdate); err != nil {
			errs = append(errs, err.Error())
		}
	}

	return errs
}

func (srv *Service) handlePartCompensate(session *schema.Session) []string {
	if len(session.Participants) == 0 {
		return nil
	}

	var errs []string

	for _, part := range session.Participants {
		if part.CompleteAction != nil {
			compensateAction := part.CompensateAction

			// if compensateAction.IsFinished() {
			// 	continue
			// }

			// @TODO: validate action

			// @TODO: update to processing

			err := srv.invokePartAction(compensateAction)
			if err != nil {
				compensateAction.Status = schema.PartActionFailed
				part.State = schema.ParticipantCompensateFailed
				errs = append(errs, err.Error())
			} else {
				compensateAction.Status = schema.PartActionCompleted
				part.State = schema.ParticipantCompensated
			}

			if _, err = srv.s.Participant().UpdateBySessionAndId(session.Id, part.Id, &schema.ParticipantUpdate{
				State:            &part.State,
				CompensateAction: part.CompensateAction,
			}); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	return errs
}
