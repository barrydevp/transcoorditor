package exclusive

import (
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type lockTableRepo struct {
	*baseRepo
	s store.LockTable
}

func NewLockTable(s store.LockTable) store.LockTable {
	return &lockTableRepo{
		baseRepo: newBaseRepo(),
		s:        s,
	}
}

func (s *lockTableRepo) Save(lockEnt *schema.LockEntry) (err error) {
	s.withLock(lockEnt.Key, func() {
		err = s.s.Save(lockEnt)
	})

	return
}

func (s *lockTableRepo) Update(lockEnt *schema.LockEntry) (err error) {
	s.withLock(lockEnt.Key, func() {
		err = s.s.Update(lockEnt)
	})
	return
}

func (s *lockTableRepo) Find(key string) (*schema.LockEntry, error) {
	return s.s.Find(key)
}

func (s *lockTableRepo) FindWithOwner(key string, owner string) (*schema.LockEntry, error) {
	return s.s.FindWithOwner(key, owner)
}

func (s *lockTableRepo) Delete(lockEnt *schema.LockEntry) (err error) {
	s.withLock(lockEnt.Key, func() {
		err = s.s.Delete(lockEnt)
	})
	return
}

func (s *lockTableRepo) DeleteByOwner(owner string) (count int64, err error) {
	// FIXME: seperated lock with key and lock with owner
	s.withLock(owner, func() {
		count, err = s.s.DeleteByOwner(owner)
	})
	return
}
