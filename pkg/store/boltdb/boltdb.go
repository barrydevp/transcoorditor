package boltdb

import (
	// "context"
	"syscall"
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

type boltdbStore struct {
	db *bbolt.DB

	session     *sessionRepo
	participant *participantRepo
}

const (
	BoltOpenTimeout        = 3
	MongoPingTimeout       = 2
	MongoDisconnectTimeout = 3
)

func newBolt() (*bbolt.DB, error) {
	bopts := &bbolt.Options{
		MmapFlags:      syscall.MAP_POPULATE,
		NoFreelistSync: true,
		Timeout:        BoltOpenTimeout * time.Second,
	}

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

	return &boltdbStore{
		db:          db,
		session:     NewSession(baseRepo),
		participant: NewParticipant(baseRepo),
	}, nil
}

func (s *boltdbStore) Close() {
	if s.db == nil {
		return
	}

	if err := s.db.Close(); err != nil {
		logger.Error("Oops... Cannot disconnect BoltDB! Reason: %v", err)

		return
	}

	logger.Info("BoltDB is disconnected.")
}

func (s *boltdbStore) Session() store.Session {
	return s.session
}

func (s *boltdbStore) Participant() store.Participant {
	return s.participant
}
