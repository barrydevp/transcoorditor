package schema

import (
	"time"
)

type ActionStatus string

const (
	ActionPending    ActionStatus = "Pending"
	ActionProcessing              = "Processing"
	ActionCompleted               = "Completed"
	ActionFailed                  = "Failed"
	ActionSkipped                 = "Skipped"
)

type ParticipantAction struct {
	Data         interface{}  `json:"data" bson:"data"`
	Target       string       `json:"target" bson:"target"`
	Status       ActionStatus `json:"status" bson:"status"`
	Result       interface{}  `json:"result" bson:"result"`
	InvokedCount int          `json:"invokedCount" bson:"invokedCount"`

	// TODO: capture invoked events
}

type ParticipantState string

const (
	ParticipantActive           ParticipantState = "Active"
	ParticipantCompensating                      = "Compensating"
	ParticipantCompensated                       = "Compensated"
	ParticipantCompensateFailed                  = "CompensateFailed"
	ParticipantCompleting                        = "Completing"
	ParticipantCompleted                         = "Completed"
	ParticipantCompleteFailed                    = "CompleteFailed"
)

type Participant struct {
	Id        string `json:"id"`
	SessionId string `json:"sessionId"`

	ClientId         string             `json:"clientId" bson:"clientId"`
	RequestId        string             `json:"requestId" bson:"requestId"`
	State            ParticipantState   `json:"state" bson:"state"`
	CompensateAction *ParticipantAction `json:"compensateAction" bson:"compensateAction"`
	CompleteAction   *ParticipantAction `json:"completeAction" bson:"completeAction"`
	UpdatedAt        *time.Time         `json:"updatedAt,omitempty" json:"bson,omitempty"`
	CreatedAt        *time.Time         `json:"createdAt" bson:"createdAt"`
}
