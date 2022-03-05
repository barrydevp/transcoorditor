package schema

import (
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/google/uuid"
)

type PartActionStatus string

const (
	PartActionPending    PartActionStatus = "Pending"
	PartActionProcessing                  = "Processing"
	PartActionCompleted                   = "Completed"
	PartActionFailed                      = "Failed"
	PartActionSkipped                     = "Skipped"
)

type PartActionResult struct {
	Error      string        `json:"error" bson:"error,omitempty"`
	StatusCode int           `json:"statusCode" bson:"statusCode,omitempty"`
	Status     string        `json:"status" bson:"status,omitempty"`
	Proto      string        `json:"proto" bson:"proto,omitempty"`
	Time       time.Duration `json:"time" bson:"time,omitempty"`
	ReceivedAt time.Time     `json:"receivedAt" bson:"receivedAt,omitempty"`
	Body       string        `json:"body" bson:"body,omitempty"`
}

type ParticipantAction struct {
	Data         interface{}         `json:"data" bson:"data"`
	Target       *string             `json:"target" bson:"target" validate:"required"`
	Status       PartActionStatus    `json:"status" bson:"status"`
	Results      []*PartActionResult `json:"results" bson:"results"`
	InvokedCount int                 `json:"invokedCount" bson:"invokedCount"`

	// TODO: capture invoked events
}

type ParticipantState string

const (
	ParticipantActive           ParticipantState = "Active"
	ParticipantCommitted                         = "Committed"
	ParticipantCompensating                      = "Compensating"
	ParticipantCompensated                       = "Compensated"
	ParticipantCompensateFailed                  = "CompensateFailed"
	ParticipantCompleting                        = "Completing"
	ParticipantCompleted                         = "Completed"
	ParticipantCompleteFailed                    = "CompleteFailed"
)

type Participant struct {
	Id        string `json:"id" bson:"id"`
	SessionId string `json:"sessionId" bson:"sessionId"`

	ClientId         string             `json:"clientId" bson:"clientId"`
	RequestId        string             `json:"requestId" bson:"requestId"`
	State            ParticipantState   `json:"state" bson:"state"`
	CompensateAction *ParticipantAction `json:"compensateAction,omitempty" bson:"compensateAction,omitempty"`
	CompleteAction   *ParticipantAction `json:"completeAction,omitempty" bson:"completeAction,omitempty"`
	UpdatedAt        *time.Time         `json:"updatedAt,omitempty" json:"bson,omitempty"`
	CreatedAt        *time.Time         `json:"createdAt" bson:"createdAt"`
}

func NewParticipant() *Participant {
	now := time.Now()

	return &Participant{
		Id:        uuid.NewString(),
		State:     ParticipantActive,
		CreatedAt: &now,
	}
}

type ParticipantJoinBody struct {
	ClientId  string `json:"clientId" validate:"required"`
	RequestId string `json:"requestId"`
}

func (p *ParticipantJoinBody) Validate() error {
	return common.GetValidate().Struct(p)
}

type ParticipantUpdate struct {
	// ClientId         *string             `json:"clientId"`
	// RequestId        *string             `json:"requestId"`
	State            *ParticipantState  `json:"state"`
	CompensateAction *ParticipantAction `json:"compensateAction"`
	CompleteAction   *ParticipantAction `json:"completeAction"`
	UpdatedAt        *time.Time         `json:"updatedAt"`
}

type ParticipantCommit struct {
	Id         *string            `json:"participantId" validate:"required"`
	Compensate *ParticipantAction `json:"compensate"`
	Complete   *ParticipantAction `json:"complete"`
}
