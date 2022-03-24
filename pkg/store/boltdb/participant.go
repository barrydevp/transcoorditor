package boltdb

import (
	// "context"
	"encoding/json"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"go.etcd.io/bbolt"
)

type participantRepo struct {
	*baseRepo
	col *collection
}

func NewParticipant(opts *baseRepo) *participantRepo {
	col := newCollection(opts.db, "session")

	return &participantRepo{
		baseRepo: opts,
		col:      col,
	}
}

func normalizeParticipant(p *schema.Participant) {
	// if p.CompensateAction != nil {
	// 	p.CompensateAction.Data = TryConvertBsonDToM(p.CompensateAction.Data)
	// }
	//
	// if p.CompleteAction != nil {
	// 	p.CompleteAction.Data = TryConvertBsonDToM(p.CompleteAction.Data)
	// }
}

func (s *participantRepo) Save(part *schema.Participant) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(part.SessionId))
		if buf == nil {
			return store.ErrSessionNotFound
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		session.Participants = append(session.Participants, part)

		buf, err = json.Marshal(session)
		if err != nil {
			return err
		}

		return b.Put([]byte(part.SessionId), buf)
	})
}

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, schemaUpdate *schema.Participant) (*schema.Participant, error) {
	var doc *schema.Participant

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(sessionId))
		if buf == nil {
			return store.ErrSessionNotFound
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		doc = session.GetParticipantAt(id)
		if doc == nil {
			return nil
		}

		needUpdate := false

		if schemaUpdate.State != "" {
			needUpdate = true
			doc.State = schemaUpdate.State
		}

		if schemaUpdate.UpdatedAt != nil {
			needUpdate = true
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		if schemaUpdate.ClientId != "" {
			needUpdate = true
			doc.ClientId = schemaUpdate.ClientId
		}

		if schemaUpdate.RequestId != "" {
			needUpdate = true
			doc.RequestId = schemaUpdate.RequestId
		}

		if schemaUpdate.CompensateAction != nil {
			needUpdate = true
			doc.CompensateAction = schemaUpdate.CompensateAction
		}

		if schemaUpdate.CompleteAction != nil {
			needUpdate = true
			doc.CompleteAction = schemaUpdate.CompleteAction
		}

		// no changes
		if !needUpdate {
			return nil
		}

		if schemaUpdate.UpdatedAt == nil {
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		buf, err = json.Marshal(session)
		if err != nil {
			return err
		}

		return b.Put([]byte(sessionId), buf)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error) {
	var doc *schema.Participant

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(sessionId))
		if buf == nil {
			return store.ErrSessionNotFound
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		doc = session.GetParticipantAt(id)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *participantRepo) FindBySessionId(sessionId string) ([]*schema.Participant, error) {
	var results []*schema.Participant

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(sessionId))
		if buf == nil {
			return store.ErrSessionNotFound
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		results = session.Participants

		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *participantRepo) FindDupInSession(sessionId string, reqPart *schema.Participant) (*schema.Participant, error) {
	allPart, err := s.FindBySessionId(sessionId)
	if err != nil {
		return nil, err
	}

	for _, part := range allPart {
		if part.RequestId == reqPart.RequestId && part.ClientId == reqPart.ClientId {
			return part, nil
		}
	}

	return nil, nil
}

func (s *participantRepo) UpdateBySessionAndId(sessionId string, id int64, schemaUpdate *schema.ParticipantUpdate) (*schema.Participant, error) {
	var doc *schema.Participant

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(sessionId))
		if buf == nil {
			return store.ErrSessionNotFound
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		doc = session.GetParticipantAt(id)
		if doc == nil {
			return nil
		}

		needUpdate := false

		if schemaUpdate.State != nil {
			needUpdate = true
			doc.State = *schemaUpdate.State
		}

		if schemaUpdate.UpdatedAt != nil {
			needUpdate = true
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		if schemaUpdate.CompensateAction != nil {
			needUpdate = true
			doc.CompensateAction = schemaUpdate.CompensateAction
		}

		if schemaUpdate.CompleteAction != nil {
			needUpdate = true
			doc.CompleteAction = schemaUpdate.CompleteAction
		}

		// no changes
		if !needUpdate {
			return nil
		}

		if schemaUpdate.UpdatedAt == nil {
			doc.UpdatedAt = schemaUpdate.UpdatedAt
		}

		buf, err = json.Marshal(session)
		if err != nil {
			return err
		}

		return b.Put([]byte(sessionId), buf)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *participantRepo) CountBySessionId(sessionId string) (int64, error) {
	allPart, err := s.FindBySessionId(sessionId)
	if err != nil {
		return 0, err
	}

	return int64(len(allPart)), nil
}

func (s *participantRepo) DeleteBySessionId(sessionId string) (int64, error) {
	deletedCount := 0

	err := s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("session"))

		buf := b.Get([]byte(sessionId))
		if buf == nil {
			// already delete
			return nil
		}

		session := &schema.Session{}
		err := json.Unmarshal(buf, session)
		if err != nil {
			return err
		}

		deletedCount = len(session.Participants)
		session.Participants = nil

		buf, err = json.Marshal(session)
		if err != nil {
			return err
		}

		return b.Put([]byte(sessionId), buf)
	})

	if err != nil {
		return 0, err
	}

	return int64(deletedCount), nil
}
