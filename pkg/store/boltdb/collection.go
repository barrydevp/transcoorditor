package boltdb

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type collection struct {
	db   *bbolt.DB
	name string
}

func newCollection(db *bbolt.DB, name string) *collection {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("cannot create bucket %s", name))
	}

	return &collection{
		db:   db,
		name: name,
	}
}
