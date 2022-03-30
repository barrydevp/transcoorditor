package replset

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type sessionRepo struct {
	*replsetBackend
	s         store.Session
	namespace string
}

func NewSession(b *replsetBackend) *sessionRepo {
	return &sessionRepo{
		replsetBackend: b,
		s:              b.s.Session(),
		namespace:      "Session",
	}
}

func (s *sessionRepo) executeRPC(c *cluster.Command) *cluster.ApplyResponse {
	method := string(c.K)

	switch method {
	case "Save":
		return s.applySave(c)
	case "PutById":
		return s.applyPutById(c)
	case "UpdateById":
		return s.applyUpdateById(c)
	case "DeleteById":
		return s.applyDeleteById(c)
	}

	return NewApplyErr(ErrRpcUnsupported)
}

func (s *sessionRepo) Save(session *schema.Session) error {
	cmd, err := cluster.NewRpcCmd(s.namespace, "Save", session)
	if err != nil {
		return err
	}

	_, err = s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRepo) applySave(c *cluster.Command) *cluster.ApplyResponse {
	session := &schema.Session{}
	err := cluster.ParseRpcCmd(c, session)
	if err != nil {
		return NewApplyErr(err)
	}

	err = s.s.Save(session)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{}
}

func (s *sessionRepo) PutById(id string, update *schema.Session) (*schema.Session, error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "PutById", id, update)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return nil, err
	}
	if doc, ok := res.(*schema.Session); ok {
		return doc, nil
	}

	return nil, ErrUnExpectedResponse
}

func (s *sessionRepo) applyPutById(c *cluster.Command) *cluster.ApplyResponse {
	id := ""
	update := &schema.Session{}
	err := cluster.ParseRpcCmd(c, &id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.PutById(id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}

func (s *sessionRepo) Find(search *schema.SessionSearch) (sessions []*schema.Session, err error) {
	return s.s.Find(search)
}

func (s *sessionRepo) FindAllUnfinished() ([]*schema.Session, error) {
	return s.s.FindAllUnfinished()
}

func (s *sessionRepo) FindById(id string) (session *schema.Session, err error) {
	session, err = s.s.FindById(id)

	return
}

func (s *sessionRepo) UpdateById(id string, update *schema.SessionUpdate) (session *schema.Session, err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "UpdateById", id, update)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return nil, err
	}
	if doc, ok := res.(*schema.Session); ok {
		return doc, nil
	}

	return nil, ErrUnExpectedResponse
}

func (s *sessionRepo) applyUpdateById(c *cluster.Command) *cluster.ApplyResponse {
	id := ""
	update := &schema.SessionUpdate{}
	err := cluster.ParseRpcCmd(c, &id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.UpdateById(id, update)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}

func (s *sessionRepo) DeleteById(id string) (session *schema.Session, err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "DeleteById", id)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return nil, err
	}
	if doc, ok := res.(*schema.Session); ok {
		return doc, nil
	}

	return nil, ErrUnExpectedResponse
}

func (s *sessionRepo) applyDeleteById(c *cluster.Command) *cluster.ApplyResponse {
	id := ""
	err := cluster.ParseRpcCmd(c, &id)
	if err != nil {
		return NewApplyErr(err)
	}

	doc, err := s.s.DeleteById(id)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: doc,
	}
}
