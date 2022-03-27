package boltdb

import (
	"encoding/json"
	"fmt"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type sessionRepo struct {
	*baseRepo
	name string
}

func NewSession(b *baseRepo) store.Session {
	name := "session"
	err := b.initCollection(name)
	if err != nil {
		panic(fmt.Sprintf("cannot create bucket %s: %v", name, err))
	}

	return &sessionRepo{
		baseRepo: b,
		name:     name,
	}
}

func (s *sessionRepo) Save(session *schema.Session) error {
	return s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		clone := *session

		return col.Put(clone.Id, clone)
	})
}

func (s *sessionRepo) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	var doc *schema.Session

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.Session{}
		_doc, err := col.Get(id, doc)
		if err != nil || _doc == nil {
			doc = nil
			return err
		}

		needUpdate := false

		if schemaUpdate.State != "" {
			needUpdate = true
			doc.State = schemaUpdate.State
		}

		if schemaUpdate.Errors != nil {
			needUpdate = true
			doc.Errors = schemaUpdate.Errors
		}

		if schemaUpdate.UpdatedAt != nil {
			needUpdate = true
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		if schemaUpdate.Timeout != 0 {
			needUpdate = true
			doc.Timeout = schemaUpdate.Timeout
		}

		if schemaUpdate.Retries != 0 {
			needUpdate = true
			doc.Retries = schemaUpdate.Retries
		}

		// no changes
		if !needUpdate {
			return nil
		}

		if schemaUpdate.UpdatedAt == nil {
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		return col.Put(id, doc)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *sessionRepo) FindById(id string) (*schema.Session, error) {
	var doc *schema.Session

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.Session{}
		_doc, err := col.Get(id, doc)
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

func (s *sessionRepo) Find(search *schema.SessionSearch) ([]*schema.Session, error) {
	var results []*schema.Session

	err := s.read(func(tx *txn) error {
		// Assume bucket exists and has keys
		col := tx.collection(s.name)

		c := col.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			doc := &schema.Session{}
			err := json.Unmarshal(v, doc)
			if err != nil {
				return err
			}
			results = append(results, doc)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *sessionRepo) FindAllUnfinished() ([]*schema.Session, error) {
	var results []*schema.Session

	err := s.read(func(tx *txn) error {
		// Assume bucket exists and has keys
		col := tx.collection(s.name)

		c := col.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			doc := &schema.Session{}
			err := json.Unmarshal(v, doc)
			if err != nil {
				return err
			}

			if !doc.IsFinished() {
				results = append(results, doc)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *sessionRepo) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (*schema.Session, error) {
	var doc *schema.Session

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.Session{}
		_doc, err := col.Get(id, doc)
		if err != nil || _doc == nil {
			doc = nil
			return err
		}

		needUpdate := false

		if schemaUpdate.State != nil {
			needUpdate = true
			doc.State = *schemaUpdate.State
		}

		if schemaUpdate.Errors != nil {
			needUpdate = true
			doc.Errors = *schemaUpdate.Errors
		}

		if schemaUpdate.UpdatedAt != nil {
			needUpdate = true
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		if schemaUpdate.Timeout != nil {
			needUpdate = true
			doc.Timeout = *schemaUpdate.Timeout
		}

		if schemaUpdate.Retries != nil {
			needUpdate = true
			doc.Retries = *schemaUpdate.Retries
		}

		// no changes
		if !needUpdate {
			return nil
		}

		if schemaUpdate.UpdatedAt == nil {
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		return col.Put(id, doc)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *sessionRepo) DeleteById(id string) (*schema.Session, error) {
	var doc *schema.Session

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &schema.Session{}
		_doc, err := col.Get(id, doc)
		if err != nil || _doc == nil {
			doc = nil
			// return err
		}

		if err := col.Delete(id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}
