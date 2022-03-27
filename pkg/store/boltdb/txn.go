package boltdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
)

type txn struct {
	tx *bbolt.Tx
}

func (t *txn) collection(name string) *collection {
	b := t.tx.Bucket([]byte(name))

	// if b == nil {
	// 	return nil, ErrCollectionNotFound
	// }

	return &collection{
		b: b,
	}
}

type collection struct {
	b *bbolt.Bucket
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
		// db:   db,
		// name: name,
	}
}

func (c *collection) Get(key string, dst interface{}) (interface{}, error) {
	docBuf := c.b.Get([]byte(key))

	if docBuf == nil {
		return nil, nil
	}

	if dst == nil {
		return docBuf, nil
	}

	err := json.Unmarshal(docBuf, dst)
	if err != nil {
		return nil, err
	}

	return docBuf, nil
}

func (c *collection) Put(key string, src interface{}) error {
	docBuf, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return c.b.Put([]byte(key), docBuf)
}

func (c *collection) Delete(key string) error {
	return c.b.Delete([]byte(key))
}

func (c *collection) Cursor() *bbolt.Cursor {
	return c.b.Cursor()
}
