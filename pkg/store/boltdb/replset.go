package boltdb

import (
	"fmt"

	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/hashicorp/raft"
)

const (
	LastLogKey = "replset_last_log"
)

type replsetRepo struct {
	*baseRepo
	name string
}

func NewReplset(b *baseRepo) store.Replset {
	name := "replset"
	err := b.initCollection(name)
	if err != nil {
		panic(fmt.Sprintf("cannot create bucket %s: %v", name, err))
	}

	return &replsetRepo{
		baseRepo: b,
		name:     name,
	}
}

func (s *replsetRepo) SaveLastLog(log *raft.Log) error {
	return s.exec(func(tx *txn) error {
		col := tx.collection(s.name)

		clone := *log

		return col.Put(LastLogKey, clone)
	})
}

func (s *replsetRepo) GetLastLog() (*raft.Log, error) {
	var doc *raft.Log

	err := s.read(func(tx *txn) error {
		col := tx.collection(s.name)

		doc = &raft.Log{}
		_doc, err := col.Get(LastLogKey, doc)
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
