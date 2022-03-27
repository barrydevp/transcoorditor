package schema

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/go-resty/resty/v2"
	// "github.com/google/uuid"
)

var (
	ErrActionRequestFailed = errors.New("action request failed")
	ErrInvalidActionUri    = fmt.Errorf("invalid action's uri. %w", exception.ErrInvalidArgument)
)

type PartActionStatus string

const (
	PartActionCreated    PartActionStatus = "Created"
	PartActionProcessing PartActionStatus = "Processing"
	PartActionCompleted  PartActionStatus = "Completed"
	PartActionFailed     PartActionStatus = "Failed"

	MAX_ACTION_INVOKED = 10
)

type PartActionResult struct {
	Error      string    `json:"error" bson:"error,omitempty"`
	StatusCode int       `json:"statusCode" bson:"statusCode,omitempty"`
	Status     string    `json:"status" bson:"status,omitempty"`
	Proto      string    `json:"proto" bson:"proto,omitempty"`
	Time       int64     `json:"time" bson:"time,omitempty"`
	ReceivedAt time.Time `json:"receivedAt" bson:"receivedAt,omitempty"`
	Body       string    `json:"body" bson:"body,omitempty"`
}

func (pr *PartActionResult) SetError(err error) {
	pr.Error = err.Error()
}

func (pr *PartActionResult) ParseRestyResp(resp *resty.Response, err error) error {
	if pr != nil {
		pr.StatusCode = resp.StatusCode()
		pr.Status = resp.Status()
		pr.Proto = resp.Proto()
		pr.Time = resp.Time().Milliseconds()
		pr.ReceivedAt = resp.ReceivedAt()
		pr.Body = resp.String()
		if resp.StatusCode() != 200 {
			return ErrActionRequestFailed
		}
	}

	return err
}

type ParticipantAction struct {
	Data         interface{}         `json:"data" bson:"data"`
	Uri          *string             `json:"uri" bson:"uri" validate:"required"`
	Status       PartActionStatus    `json:"status" bson:"status"`
	Results      []*PartActionResult `json:"results" bson:"results"`
	InvokedCount int                 `json:"invokedCount" bson:"invokedCount"`

	// TODO: capture invoked events
}

func (pa *ParticipantAction) IsFinished() bool {
	return pa.Status == PartActionCompleted || pa.InvokedCount > MAX_ACTION_INVOKED
}

func (pa *ParticipantAction) requestActionHTTP() (*resty.Response, error) {
	// build request
	req := util.GetRequest().R()

	if pa.Data != nil {
		req.SetBody(pa.Data)
	}

	return req.Post(*pa.Uri)
}

func (pa *ParticipantAction) ValidateAction() error {
	if pa.Uri == nil {
		return ErrInvalidActionUri
	}

	if _, err := url.ParseRequestURI(*pa.Uri); err != nil {
		return ErrInvalidActionUri
	}

	return nil
}

// invoke participant action and update it's result
func (pa *ParticipantAction) InvokePartAction() error {
	result := &PartActionResult{}

	if pa.Status == PartActionCompleted {
		return nil
	}

	// validate action
	err := pa.ValidateAction()

	if err == nil {
		err = result.ParseRestyResp(pa.requestActionHTTP())
	}

	if err != nil {
		pa.Status = PartActionFailed
		result.SetError(err)
	} else {
		pa.Status = PartActionCompleted
	}

	pa.Results = append(pa.Results, result)
	pa.InvokedCount++

	return err
}

type ParticipantState string

const (
	ParticipantActive           ParticipantState = "Active"
	ParticipantCommitted        ParticipantState = "Committed"
	ParticipantCompensating     ParticipantState = "Compensating"
	ParticipantCompensated      ParticipantState = "Compensated"
	ParticipantCompensateFailed ParticipantState = "CompensateFailed"
	ParticipantCompleting       ParticipantState = "Completing"
	ParticipantCompleted        ParticipantState = "Completed"
	ParticipantCompleteFailed   ParticipantState = "CompleteFailed"
)

type Participant struct {
	Id        int64  `json:"id" bson:"id"`
	SessionId string `json:"sessionId" bson:"sessionId"`

	ClientId         string             `json:"clientId" bson:"clientId"`
	RequestId        string             `json:"requestId" bson:"requestId"`
	State            ParticipantState   `json:"state" bson:"state"`
	CompensateAction *ParticipantAction `json:"compensateAction,omitempty" bson:"compensateAction,omitempty"`
	CompleteAction   *ParticipantAction `json:"completeAction,omitempty" bson:"completeAction,omitempty"`
	UpdatedAt        *time.Time         `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedAt        *time.Time         `json:"createdAt" bson:"createdAt"`
}

func NewParticipant() *Participant {
	now := time.Now()

	return &Participant{
		// Id:        uuid.NewString(),
		Id:        0,
		State:     ParticipantActive,
		CreatedAt: &now,
	}
}

func (p *Participant) GetAction(compensate bool) *ParticipantAction {
	if compensate {
		return p.CompensateAction
	}

	return p.CompleteAction
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
	Id         *int64             `json:"participantId" validate:"required"`
	Compensate *ParticipantAction `json:"compensate"`
	Complete   *ParticipantAction `json:"complete"`
}
