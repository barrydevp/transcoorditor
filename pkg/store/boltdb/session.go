package boltdb

import (
	"encoding/json"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"go.etcd.io/bbolt"
)

type sessionRepo struct {
	*baseRepo
	col *collection
}

func NewSession(opts *baseRepo) *sessionRepo {
	col := newCollection(opts.db, "session")

	return &sessionRepo{
		baseRepo: opts,
		col:      col,
	}
}

func (s *sessionRepo) Save(session *schema.Session) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		clone := *session
		clone.Participants = nil

		buf, err := json.Marshal(&clone)
		if err != nil {
			return err
		}

		return b.Put([]byte(clone.Id), buf)
	})
}

func (s *sessionRepo) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	var doc *schema.Session

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(id))

		if buf == nil {
			return nil
		}

		doc = &schema.Session{}
		err := json.Unmarshal(buf, doc)
		if err != nil {
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

		buf, err = json.Marshal(doc)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), buf)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *sessionRepo) FindById(id string) (*schema.Session, error) {
	var doc *schema.Session

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(id))
		if buf == nil {
			return nil
		}

		doc = &schema.Session{}
		err := json.Unmarshal(buf, doc)
		if err != nil {
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

	err := s.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("session"))

		c := b.Cursor()
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

	err := s.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("session"))

		c := b.Cursor()
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

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(id))

		if buf == nil {
			return nil
		}

		doc = &schema.Session{}
		err := json.Unmarshal(buf, doc)
		if err != nil {
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

		buf, err = json.Marshal(doc)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), buf)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *sessionRepo) DeleteById(id string) (*schema.Session, error) {
	var doc *schema.Session

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(id))
		if buf == nil {
			return nil
		}

		err := b.Delete([]byte(id))
		if err != nil {
			return err
		}

		doc = &schema.Session{}
		err = json.Unmarshal(buf, doc)
		if err != nil {
			// we ignore this error because we already delete the record
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}
