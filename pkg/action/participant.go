package action

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func (ac *Action) invokeAction(a *schema.Action, handler ActionHandler) (resultCh chan interface{}, err error) {
//
// 	// check if status is not pending to run
// 	if a.Status != "Pending" {
// 		return nil, util.NewError("action was %v", a.Status)
// 	}
//
// 	// increase invoked count
// 	a.InvokedCount++
//
// 	a.Status = schema.ActionProcessing
//
// 	// invoke handler
// 	resultCh = make(chan interface{})
//
// 	go func() {
// 		result, err := handler(a)
//
// 		if err != nil {
//
// 		}
// 	}()
//
// 	if err != nil {
// 		a.Status = ActionCompleted
//
// 	}
//
// 	return
// }

var (
	ErrParticipantNotFound = errors.New("participant not found")
	ErrInvalidActionPath   = errors.New("invalid action's path")
	ErrActionRequestFailed = errors.New("action request failed")
)

func (ac *Action) findParticipantById(id string) (*schema.Participant, error) {
	doc, err := ac.s.Participant().FindById(id)
	if err != nil {
		return nil, util.Errorf("failed to get participant: %w", err)
	}

	if doc == nil {
		return nil, ErrParticipantNotFound
	}

	return doc, nil
}

func (ac *Action) invokePartAction(action *schema.ParticipantAction) (err error) {
	var (
		req  *resty.Request
		resp *resty.Response
	)

	action.InvokedCount++

	result := &schema.PartActionResult{}

	if action.Target == nil {
		err = ErrInvalidActionPath
		goto RET_ERROR
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

	resp, err = req.Post(*action.Target)

	if err != nil {
		goto RET_ERROR
	}

	result.StatusCode = resp.StatusCode()
	result.Status = resp.Status()
	result.Proto = resp.Proto()
	result.Time = resp.Time()
	result.ReceivedAt = resp.ReceivedAt()
	result.Body = resp.String()

	if resp.StatusCode() != 200 {
		err = ErrActionRequestFailed
		goto RET_ERROR
	}

	return nil

RET_ERROR:
	result.Error = err.Error()
	action.Results = append(action.Results, result)

	return err
}

func (ac *Action) handlePartComplete(session *schema.Session) error {
	if len(session.Participants) == 0 {
		return nil
	}

	for _, part := range session.Participants {
		if part.CompleteAction != nil {
			completeAc := part.CompleteAction

			// @TODO: validate action

			// @TODO: update to processing

			err := ac.invokePartAction(completeAc)
			if err != nil {
				completeAc.Status = schema.PartActionFailed
				part.State = schema.ParticipantCompleteFailed
			} else {
				completeAc.Status = schema.PartActionCompleted
				part.State = schema.ParticipantCompleted
			}

			if _, err = ac.s.Participant().UpdateById(part.Id, &schema.ParticipantUpdate{
				State:          &part.State,
				CompleteAction: part.CompleteAction,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
