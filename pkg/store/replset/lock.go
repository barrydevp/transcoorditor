package replset

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type lockTableRepo struct {
	*replsetBackend
	s         store.LockTable
	namespace string
}

func NewLockTable(b *replsetBackend) *lockTableRepo {
	return &lockTableRepo{
		replsetBackend: b,
		s:              b.s.LockTable(),
		namespace:      "LockTable",
	}
}

func (s *lockTableRepo) executeRPC(c *cluster.Command) *cluster.ApplyResponse {
	method := string(c.K)

	switch method {
	case "Save":
		return s.applySave(c)
	case "Update":
		return s.applyUpdate(c)
	case "Delete":
		return s.applyDelete(c)
	case "DeleteByOwner":
		return s.applyDeleteByOwner(c)
	}

	return NewApplyErr(ErrRpcUnsupported)
}

func (s *lockTableRepo) Save(lockEnt *schema.LockEntry) error {
	cmd, err := cluster.NewRpcCmd(s.namespace, "Save", lockEnt)
	if err != nil {
		return err
	}

	_, err = s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) applySave(c *cluster.Command) *cluster.ApplyResponse {
	lockEnt := &schema.LockEntry{}
	err := cluster.ParseRpcCmd(c, lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	err = s.s.Save(lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{}
}

func (s *lockTableRepo) Update(lockEnt *schema.LockEntry) error {
	cmd, err := cluster.NewRpcCmd(s.namespace, "Update", lockEnt)
	if err != nil {
		return err
	}

	_, err = s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) applyUpdate(c *cluster.Command) *cluster.ApplyResponse {
	lockEnt := &schema.LockEntry{}
	err := cluster.ParseRpcCmd(c, lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	err = s.s.Update(lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{}
}

func (s *lockTableRepo) Find(key string) (*schema.LockEntry, error) {
	return s.s.Find(key)
}

func (s *lockTableRepo) FindWithOwner(key string, owner string) (*schema.LockEntry, error) {
	return s.s.FindWithOwner(key, owner)
}

func (s *lockTableRepo) Delete(lockEnt *schema.LockEntry) (err error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "Delete", lockEnt)
	if err != nil {
		return err
	}

	_, err = s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) applyDelete(c *cluster.Command) *cluster.ApplyResponse {
	lockEnt := &schema.LockEntry{}
	err := cluster.ParseRpcCmd(c, lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	err = s.s.Delete(lockEnt)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{}
}

func (s *lockTableRepo) DeleteByOwner(owner string) (int64, error) {
	cmd, err := cluster.NewRpcCmd(s.namespace, "DeleteByOnwer", owner)
	if err != nil {
		return 0, err
	}

	res, err := s.c.Execute(cmd, executeTimeout)
	if err != nil {
		return 0, err
	}
	if count, ok := res.(int64); ok {
		return count, nil
	}
	if err != nil {
		return 0, err
	}

	return 0, ErrUnExpectedResponse
}

func (s *lockTableRepo) applyDeleteByOwner(c *cluster.Command) *cluster.ApplyResponse {
	owner := ""
	err := cluster.ParseRpcCmd(c, &owner)
	if err != nil {
		return NewApplyErr(err)
	}

	count, err := s.s.DeleteByOwner(owner)
	if err != nil {
		return NewApplyErr(err)
	}

	return &cluster.ApplyResponse{
		Res: count,
	}
}
