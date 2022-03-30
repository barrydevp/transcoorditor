package replset

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type participantRepo struct {
	*replsetBackend
	s         store.Participant
	namespace string
}

func NewParticipant(b *replsetBackend) *participantRepo {
	return &participantRepo{
		replsetBackend: b,
		s:              b.s.Participant(),
		namespace:      "Participant",
	}
}

func (s *participantRepo) executeRPC(c *cluster.Command) *cluster.ApplyResponse {
	method := string(c.K)

	switch method {
	case "Save":
		return s.applySave(c)
	case "PutBySessionAndId":
		return s.applyPutBySessionAndId(c)
	case "UpdateBySessionAndId":
		return s.applyUpdateBySessionAndId(c)
	case "DeleteBySessionId":
		return s.applyDeleteBySessionId(c)
	}

	return NewApplyErr(ErrRpcUnsupported)
}

func (s *participantRepo) applySave(c *cluster.Command) *cluster.ApplyResponse {
	part := &schema.Participant{}
	err := cluster.ParseRpcCmd(c, part)
	if err != nil {
		return NewApplyErr(err)
	}

	err = s.s.Save(part)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{}
}

func (s *participantRepo) Save(part *schema.Participant) error {
	cmd, err := cluster.NewRpcCmd(s.namespace, "Save", part)
	if err != nil {
		return err
	}

	_, err = s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, update *schema.Participant) (pa *schema.Participant, err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "PutBySessionAndId", sessionId, id, update)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return nil, err
	}
	if doc, ok := res.(*schema.Participant); ok {
		return doc, nil
	}

	return nil, ErrUnExpectedResponse
}

func (s *participantRepo) applyPutBySessionAndId(c *cluster.Command) *cluster.ApplyResponse {
	sessionId := ""
	id := int64(0)
	update := &schema.Participant{}
	err := cluster.ParseRpcCmd(c, &sessionId, &id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.PutBySessionAndId(sessionId, id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (part *schema.Participant, err error) {
	part, err = s.s.FindBySessionAndId(sessionId, id)

	return
}

func (s *participantRepo) FindBySessionId(sessionId string) (parts []*schema.Participant, err error) {
	parts, err = s.s.FindBySessionId(sessionId)

	return
}

func (s *participantRepo) FindDupInSession(sessionId string, part *schema.Participant) (pa *schema.Participant, err error) {
	pa, err = s.s.FindDupInSession(sessionId, part)

	return
}

func (s *participantRepo) UpdateBySessionAndId(sessionId string, id int64, update *schema.ParticipantUpdate) (part *schema.Participant, err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "UpdateBySessionAndId", sessionId, id, update)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return nil, err
	}
	if doc, ok := res.(*schema.Participant); ok {
		return doc, nil
	}

	return nil, ErrUnExpectedResponse
}

func (s *participantRepo) applyUpdateBySessionAndId(c *cluster.Command) *cluster.ApplyResponse {
	sessionId := ""
	id := int64(0)
	update := &schema.ParticipantUpdate{}
	err := cluster.ParseRpcCmd(c, &sessionId, &id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.UpdateBySessionAndId(sessionId, id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}

func (s *participantRepo) CountBySessionId(sessionId string) (count int64, err error) {
	count, err = s.s.CountBySessionId(sessionId)

	return
}

func (s *participantRepo) DeleteBySessionId(sessionId string) (count int64, err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "DeleteBySessionId", sessionId)
	if err != nil {
		return 0, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return 0, err
	}
	if doc, ok := res.(int64); ok {
		return doc, nil
	}

	return 0, ErrUnExpectedResponse
}

func (s *participantRepo) applyDeleteBySessionId(c *cluster.Command) *cluster.ApplyResponse {
	sessionId := ""
	err := cluster.ParseRpcCmd(c, &sessionId)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.DeleteBySessionId(sessionId)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}
