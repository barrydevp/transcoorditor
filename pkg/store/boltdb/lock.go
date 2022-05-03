package boltdb

import (
	"fmt"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type lockTableRepo struct {
	*baseRepo
	name string
}

func NewLockTable(b *baseRepo) store.LockTable {
	name := "locktable"
	err := b.initCollection(name)
	if err != nil {
		panic(fmt.Sprintf("cannot create bucket %s: %v", name, err))
	}

	return &lockTableRepo{
		baseRepo: b,
		name:     name,
	}
}

func (s *lockTableRepo) Save(lockEnt *schema.LockEntry) error {
	return s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc := &schema.LockEntry{}
		_, err := col.Get(lockEnt.Key, doc)
		if err != nil {
			return err
		}

		// if _doc != nil {
		// 	if !doc.IsExpired() {
		// 		return fmt.Errorf("%w: owner(%v)", store.ErrLockExists, doc.Owner)
		// 	}
		// }

		clone := *lockEnt

		return col.Put(clone.Key, clone)
	})
}

func (s *lockTableRepo) Update(lockEnt *schema.LockEntry) error {
	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc := &schema.LockEntry{}
		_doc, err := col.Get(lockEnt.Key, doc)
		if err != nil {
			return err
		}

		if _doc == nil {
			return store.ErrLockNotFound
		}

		if lockEnt.Owner != doc.Owner {
			return store.ErrLockNotOwner
		}

		// if lockEnt.IsExpired() {
		// 	return store.ErrLockExpired
		// }

		doc.ExpiredAt = lockEnt.ExpiredAt

		return col.Put(lockEnt.Key, doc)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) Find(key string) (*schema.LockEntry, error) {
	var doc *schema.LockEntry

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.LockEntry{}
		_doc, err := col.Get(key, doc)
		if err != nil || _doc == nil {
			doc = nil
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *lockTableRepo) FindWithOwner(key string, owner string) (*schema.LockEntry, error) {
	var doc *schema.LockEntry

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.LockEntry{}
		_doc, err := col.Get(key, doc)
		if err != nil || _doc == nil {
			doc = nil
			return err
		}

		if doc.Owner != owner {
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *lockTableRepo) Delete(lockEnt *schema.LockEntry) error {
	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc := &schema.LockEntry{}
		_doc, err := col.Get(lockEnt.Key, doc)
		if err != nil || _doc == nil {
			doc = nil
			return err
		}

		if doc.Owner != lockEnt.Owner {
			return nil
			// return store.ErrLockNotOwner
		}

		if err := col.Delete(lockEnt.Key); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// TBD
func (s *lockTableRepo) DeleteByOwner(owner string) (int64, error) {
	// err := s.exec(func(tx *txn) error {
	// 	col := tx.collection(s.name)
	//
	// 	doc := &schema.LockEntry{}
	// 	_doc, err := col.Get(lockEnt.Key, doc)
	// 	if err != nil || _doc == nil {
	// 		doc = nil
	// 		return err
	// 	}
	//
	// 	if doc.Owner != lockEnt.Owner {
	// 		return nil
	// 		// return store.ErrLockNotOwner
	// 	}
	//
	// 	if err := col.Delete(lockEnt.Key); err != nil {
	// 		return err
	// 	}
	//
	// 	return nil
	// })
	//
	// if err != nil {
	// 	return err
	// }

	return 0, nil
}
