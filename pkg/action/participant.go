package action

import (
	"time"
	// "github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/google/uuid"
)

type ActionHandler func(a *schema.ParticipantAction) (interface{}, error)

func NewParticipant(clientId string) *schema.Participant {
	now := time.Now()

	return &schema.Participant{
		Id:        uuid.NewString(),
		ClientId:  clientId,
		State:     schema.ParticipantActive,
		CreatedAt: &now,
	}
}

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
