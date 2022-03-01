package action

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
)

// "time"
// "github.com/barrydevp/transcoorditor/pkg/util"
// "github.com/barrydevp/transcoorditor/pkg/schema"
// "github.com/google/uuid"

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
	ErrParticipantNotFound error = errors.New("participant not found")
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

