package boltdb

import (
	// "context"
	"fmt"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type participantRepo struct {
	*baseRepo
	name string
}

func NewParticipant(b *baseRepo) store.Participant {
	name := "session"
	err := b.initCollection(name)
	if err != nil {
		panic(fmt.Sprintf("cannot create bucket %s: %v", name, err))
	}

	return &participantRepo{
		baseRepo: b,
		name:     name,
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

func (s *participantRepo) getSession(col *collection, id string) (*schema.Session, error) {
	doc := &schema.Session{}
	_doc, err := col.Get(id, doc)
	if err != nil {
		doc = nil
		return nil, err
	}
	if _doc == nil {
		return nil, store.ErrSessionNotFound
	}

	return doc, nil
}

func (s *participantRepo) Save(part *schema.Participant) error {
	return s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		doc, err := s.getSession(col, part.SessionId)
		if err != nil {
			return err
		}

		if doc == nil {

		}

		doc.Participants = append(doc.Participants, part)

		return col.Put(doc.Id, doc)
	})
}

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, schemaUpdate *schema.Participant) (*schema.Participant, error) {
	var doc *schema.Participant

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		session, err := s.getSession(col, sessionId)
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

		return col.Put(sessionId, session)
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error) {
	var doc *schema.Participant

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		session, err := s.getSession(col, sessionId)
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

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		session, err := s.getSession(col, sessionId)
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

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		session, err := s.getSession(col, sessionId)
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

		return col.Put(sessionId, session)
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

	err := s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		session, err := s.getSession(col, sessionId)
		if err != nil {
			return nil
			// return err
		}

		deletedCount = len(session.Participants)
		session.Participants = nil

		return col.Put(sessionId, session)
	})

	if err != nil {
		return 0, err
	}

	return int64(deletedCount), nil
}
