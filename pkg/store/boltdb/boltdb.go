package boltdb

import (
	// "context"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.etcd.io/bbolt"
)

var logger = common.Logger().WithFields(logrus.Fields{
	"pkg": "store/boltdb",
})

type baseRepo struct {
	db *bbolt.DB
}

func (b *baseRepo) initCollection(name string) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (b *baseRepo) read(fn func(tx *txn) error) error {
	return b.db.View(func(tx *bbolt.Tx) error {
		return fn(&txn{
			tx: tx,
		})
	})
}

func (b *baseRepo) exec(fn func(tx *txn) error) error {
	return b.db.Batch(func(tx *bbolt.Tx) error {
		return fn(&txn{
			tx: tx,
		})
	})
}

type boltdbBackend struct {
	*store.Backend
	db *bbolt.DB
}

const (
	BoltOpenTimeout = 3
)

func newBolt() (*bbolt.DB, error) {
	// default options, different on each arch and os
	bopts := boltOpenOptions

	bopts.Timeout = BoltOpenTimeout * time.Second

	// https://github.com/etcd-io/etcd/blob/4504daa6b03f361a44a760e4bd5caf0e4ae9fc7e/server/storage/backend/backend.go#L41
	bopts.InitialMmapSize = 10 * 1024 * 1024 * 1024
	// bopts.FreelistType = bbolt.FreelistArrayType
	// bopts.NoSync = bcfg.UnsafeNoFsync
	// bopts.NoGrowSync = bcfg.UnsafeNoFsync
	// bopts.Mlock = bcfg.Mlock

	return bbolt.Open(viper.GetString("BOLTDB_PATH"), 0600, bopts)
}

func NewStore() (store.Interface, error) {
	db, err := newBolt()
	if err != nil {
		return nil, err
	}

	logger.Info("BoltDB is connected.")

	baseRepo := &baseRepo{
		db: db,
	}

	backend := &store.Backend{
		SessionImpl:     NewSession(baseRepo),
		ParticipantImpl: NewParticipant(baseRepo),
		ReplsetImpl:     NewReplset(baseRepo),
		LockTableImpl:   NewLockTable(baseRepo),
	}

	return &boltdbBackend{
		Backend: backend,
		db:      db,
	}, nil
}

func (s *boltdbBackend) Close() {
	if s.db == nil {
		return
	}

	if err := s.db.Close(); err != nil {
		logger.Error("Oops... Cannot disconnect BoltDB! Reason: %v", err)

		return
	}

	logger.Info("BoltDB is disconnected.")
}
