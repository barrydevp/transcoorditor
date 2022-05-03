package memory

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
)

// memory storage
// TBD
type lockTableRepo struct {
	m map[string]interface{}
}

func NewLockTable() *lockTableRepo {

	return &lockTableRepo{
		m: make(map[string]interface{}),
	}
}

func (s *lockTableRepo) Save(lockEnt *schema.LockEntry) error {
	return nil
}

func (s *lockTableRepo) Update(lockEnt *schema.LockEntry) error {
	return nil
}

func (s *lockTableRepo) Find(key string) (*schema.LockEntry, error) {
	return nil, nil
}

func (s *lockTableRepo) FindWithOwner(key string, owner string) (*schema.LockEntry, error) {
	return nil, nil
}

func (s *lockTableRepo) Delete(lockEnt *schema.LockEntry) error {
	return nil
}

func (s *lockTableRepo) DeleteByOwner(owner string) (int64, error) {
	return 0, nil
}
